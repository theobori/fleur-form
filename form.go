package fleurform

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
	gserver "github.com/theobori/fleur/gopher/server"
	"github.com/theobori/fleur/gophermap"
	"github.com/theobori/fleur/server"
)

type SubmitCallback func(Parameters, *server.Server, *server.RequestContext) error

type Form struct {
	prefixPath      string
	parametersNames []string
	submitCallback  SubmitCallback
	submitName      string
}

func NewForm(prefixPath string, parametersNames []string, submitCallback SubmitCallback) *Form {
	prefixPath = "/" + strings.Trim(prefixPath, "/") + "/"

	return &Form{
		prefixPath:      prefixPath,
		parametersNames: parametersNames,
		submitCallback:  submitCallback,
		submitName:      uuid.NewString(),
	}
}

func (f *Form) FormPath() string {
	return f.prefixPath + "form"
}

func (f *Form) formCallback(server *server.Server, ctx *server.RequestContext) error {
	var (
		parameters           Parameters
		parametersString     string
		parametersNamesuffix string
	)

	virtualPathSplitted := strings.Split(ctx.VirtualPath, ParametersBegin)
	if len(virtualPathSplitted) >= 2 {
		parametersString = virtualPathSplitted[1]
		parametersNamesuffix = ParametersBegin + parametersString
		parameters = GetParametersFromString(parametersString)
	} else {
		parametersString = ""
		parametersNamesuffix = ""
		parameters = Parameters{}
	}

	items := []*gophermap.Item{}

	for _, parameterName := range f.parametersNames {
		var item *gophermap.Item

		value, ok := parameters[parameterName]
		if !ok {
			item = server.NewItem(
				gophermap.ItemTypeGopherFullTextSearch,
				parameterName+":",
				f.prefixPath+parameterName+parametersNamesuffix,
			)
		} else {
			item = server.NewItem(
				gophermap.ItemTypeInlineText,
				parameterName+": "+value,
				f.FormPath(),
			)
		}

		items = append(items, item)
	}

	items = append(
		items,
		server.NewItem(
			gophermap.ItemTypeInlineText,
			"",
			"/",
		),
	)
	items = append(
		items,
		server.NewItem(
			gophermap.ItemTypeGopherMenu,
			"Submit",
			f.prefixPath+f.submitName+ParametersBegin+parametersString,
		),
	)

	menu := gophermap.RenderMenu(items...)

	return gserver.SendString(ctx.Conn, menu)
}

func (f *Form) setParametersRoutes(router *server.Router) error {
	for _, keyword := range f.parametersNames {
		err := router.SetWithWeight(
			10,
			"^"+f.prefixPath+keyword+".*",
			func(server *server.Server, ctx *server.RequestContext) error {
				if len(ctx.SearchParameter) == 0 {
					return server.SendError(ctx.Conn, "Missing a search parameter.")
				}

				hasParametersBegin := strings.Contains(ctx.VirtualPath, ParametersBegin)

				var preKeyword string
				if !hasParametersBegin {
					preKeyword = ParametersBegin
				} else {
					preKeyword = ParametersSeparator
				}

				ctx.VirtualPath = fmt.Sprintf(
					"%s%s%s%s%s",
					ctx.VirtualPath,
					preKeyword,
					keyword,
					PairSeparator,
					ctx.SearchParameter,
				)

				return f.formCallback(server, ctx)
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}

func (f *Form) setFormRoute(router *server.Router) error {
	return router.SetWithWeight(
		1,
		f.FormPath()+".*",
		f.formCallback,
	)
}

func (f *Form) setSubmitRoute(router *server.Router) error {
	return router.SetWithWeight(
		1,
		"^"+f.prefixPath+f.submitName+".*",
		func(server *server.Server, ctx *server.RequestContext) error {
			parameters := GetParametersFromPath(ctx.VirtualPath)

			return f.submitCallback(parameters, server, ctx)
		},
	)
}

func (f *Form) SetRoutes(router *server.Router) error {
	var err error

	err = f.setParametersRoutes(router)
	if err != nil {
		return err
	}

	err = f.setFormRoute(router)
	if err != nil {
		return err
	}

	err = f.setSubmitRoute(router)
	if err != nil {
		return err
	}

	return nil
}

func AddFormToRouter(router *server.Router, prefixPath string, parametersNames []string, submitCallback SubmitCallback) error {
	form := NewForm(prefixPath, parametersNames, submitCallback)

	err := form.SetRoutes(router)
	if err != nil {
		return err
	}

	return nil
}
