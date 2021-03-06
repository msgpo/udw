package udwRpcGoAndObjectiveCBuilder

import (
	"bytes"
	"github.com/tachyon-protocol/udw/udwCache"
	"github.com/tachyon-protocol/udw/udwFile"
	"github.com/tachyon-protocol/udw/udwGoSource/udwGoParser"
	"github.com/tachyon-protocol/udw/udwJson"
	"github.com/tachyon-protocol/udw/udwRpc/udwRpcSameProcess/udwRspBuilderLib"
	"path/filepath"
	"strings"
)

type MustBuildRequest struct {
	PackagePath string

	Filter func(name string) bool `json:"-"`

	GoToOcFunctionList []string

	OutFileName string

	BuildFlagContent string
}

func MustBuild(req MustBuildRequest) {
	udwCache.MustMd5FileChangeCache("udwRpcGoAndObjectiveC_"+udwJson.MustMarshalToString(req), []string{
		filepath.Join("src", req.PackagePath),
		"src/github.com/tachyon-protocol/udw/udwRpc/udwRpcGoAndObjectiveC",
	}, func() {
		MustBuildNoCache(req)
	})
}

func filterTrue(name string) bool {
	return true
}

func MustBuildNoCache(req MustBuildRequest) {
	if req.Filter == nil {
		req.Filter = filterTrue
	}
	if req.OutFileName == "" {
		req.OutFileName = "zzzig_udwRpcGoAndObjectiveC"
	}
	if req.BuildFlagContent == "" {
		req.BuildFlagContent = "darwin"
	}
	filter := func(name string) bool {
		if strings.HasPrefix(name, "udwGaoc_go_") {
			return false
		}
		if strings.HasPrefix(name, "_udwRsp_") {
			return false
		}
		return req.Filter(name)
	}
	pkg := udwGoParser.MustParsePackegeFromDefaultEnv(req.PackagePath)

	ctx := builderCtx{
		req:                  req,
		seenMarhshalNameMap:  map[string]bool{},
		seenUnmarshalNameMap: map[string]bool{},
	}
	outDir := pkg.GetDirPath()
	ctx.goBuilderCtx = udwRspBuilderLib.NewGoBuilderCtx(udwRspBuilderLib.GoBuilderCtxRequest{
		OutGoFilePath:                 filepath.Join(outDir, req.OutFileName+".go"),
		IsNoParameterFromGoDirectCall: true,
		IsStringUTF16:                 false,
		PkgImportPath:                 pkg.GetPkgImportPath(),
		BuildFlagContent:              req.BuildFlagContent,
		CgoHeaderContent: `#cgo CFLAGS: -x objective-c -fobjc-arc
#cgo LDFLAGS: -framework Foundation
#include <Foundation/Foundation.h>
`,
	})
	ctx.goBuilderCtx.GoFileContext.AddUnderScoreImportPath("github.com/tachyon-protocol/udw/udwRpc/udwRpcGoAndObjectiveC/udwGaocLib")

	ctx.hFileBuffer.WriteString(`// +build ` + req.BuildFlagContent + `
// Auto generated by udwRpcGoAndObjectiveC , do not edit it.

#include <Foundation/Foundation.h>

`)
	toGoFnVisitor := func(f func(fnDef *udwGoParser.FuncOrMethodDeclaration)) {
		ctx.goBuilderCtx.CurrentProcessFnName = ""
		for _, fnDef := range pkg.GetFuncList() {
			if filter(fnDef.GetName()) == false {
				continue
			}
			ctx.goBuilderCtx.CurrentProcessFnName = fnDef.GetName()
			f(fnDef)
		}
	}
	fromGoFnVisitor := func(f func(fnDef *udwGoParser.FuncOrMethodDeclaration)) {
		ctx.goBuilderCtx.CurrentProcessFnName = ""
		for _, fContent := range req.GoToOcFunctionList {
			fnDef := udwGoParser.MustParseGoFuncDeclaration(udwGoParser.MustParseGoFuncDeclarationRequest{
				Pkg:                      pkg,
				GoFuncDeclarationContent: fContent,
			})
			ctx.goBuilderCtx.CurrentProcessFnName = fnDef.GetName()
			f(fnDef)
		}
	}
	ctx.genGoStructToOc(func(f func(fnDef *udwGoParser.FuncOrMethodDeclaration)) {
		toGoFnVisitor(f)
		fromGoFnVisitor(f)
	})

	ctx.goBuilderCtx.GenGo(toGoFnVisitor, fromGoFnVisitor)

	toGoFnVisitor(func(fn *udwGoParser.FuncOrMethodDeclaration) {
		ctx.ocToGoGenFnHAndMFile(fn)
	})
	fromGoFnVisitor(func(fn *udwGoParser.FuncOrMethodDeclaration) {
		ctx.goToOcGenFnHAndMFile(fn)
	})
	ctx.goBuilderCtx.CurrentProcessFnName = ""

	udwFile.MustCheckContentAndWriteFileWithMkdir(filepath.Join(outDir, req.OutFileName+".h"), ctx.hFileBuffer.Bytes())
	udwFile.MustCheckContentAndWriteFileWithMkdir(filepath.Join(outDir, req.OutFileName+".m"), []byte(`// +build `+req.BuildFlagContent+`
// Auto generated by udwRpcGoAndObjectiveC , do not edit it.

#include "_cgo_export.h"
#include "src/github.com/tachyon-protocol/udw/udwRpc/udwRpcSameProcess/udwRspLib/udwRspLib.h"
#include "src/github.com/tachyon-protocol/udw/udwRpc/udwRpcGoAndObjectiveC/udwGaocLib/udwGaocLib.h"
#include "`+req.OutFileName+`.h"
`+ctx.mFuncFileBuffer.String()+"\n"+ctx.mFileBuffer.String()))
}

type builderCtx struct {
	req                  MustBuildRequest
	goBuilderCtx         *udwRspBuilderLib.GoBuilderCtx
	hFileBuffer          bytes.Buffer
	mFileBuffer          bytes.Buffer
	mFuncFileBuffer      bytes.Buffer
	seenMarhshalNameMap  map[string]bool
	seenUnmarshalNameMap map[string]bool
}

func (ctx *builderCtx) getCurrentProcessFnName() string {
	return ctx.goBuilderCtx.CurrentProcessFnName
}

func (ctx *builderCtx) getNextVarString() string {
	return ctx.goBuilderCtx.GetNextVarString()
}
