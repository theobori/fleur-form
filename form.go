package fleurform

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	gserver "github.com/theobori/fleur/gopher/server"
	"github.com/theobori/fleur/gophermap"
	"github.com/theobori/fleur/server"
)

type SubmitCallback func(Parameters, *server.Server, *server.RequestContext) error

type Form struct {
	prefixPath         string
	parametersMetadata []ParameterMetadata
	submitCallback     SubmitCallback
	submitName         string
}

func NewForm(
	prefixPath string,
	parametersMetadata []ParameterMetadata,
	submitCallback SubmitCallback,
) *Form {
	prefixPath = "/" + strings.Trim(prefixPath, "/")

	return &Form{
		prefixPath:         prefixPath,
		parametersMetadata: parametersMetadata,
		submitCallback:     submitCallback,
		submitName:         uuid.NewString(),
	}
}

func (f *Form) FormPath() string {
	return filepath.Join(f.prefixPath, "form")
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

	for _, parameterMetadata := range f.parametersMetadata {
		var item *gophermap.Item

		value, ok := parameters[parameterMetadata.Name]
		if !ok {
			item = server.NewItem(
				gophermap.ItemTypeGopherFullTextSearch,
				parameterMetadata.String()+":",
				f.FormPath()+"/"+parameterMetadata.Name+parametersNamesuffix,
			)
		} else {
			item = server.NewItem(
				gophermap.ItemTypeInlineText,
				parameterMetadata.String()+": "+value,
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
			f.FormPath()+"/"+f.submitName+ParametersBegin+parametersString,
		),
	)

	menu := gophermap.RenderMenu(items...)

	return gserver.SendString(ctx.Conn, menu)
}

func (f *Form) setParametersRoutes(router *server.Router) error {
	for _, parameterMetadata := range f.parametersMetadata {
		err := router.SetWithWeight(
			10,
			"^"+f.FormPath()+"/"+parameterMetadata.Name+".*",
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
					parameterMetadata.Name,
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
		0,
		"^"+f.FormPath()+".*",
		f.formCallback,
	)
}

func (f *Form) setSubmitRoute(router *server.Router) error {
	return router.SetWithWeight(
		1,
		"^"+f.FormPath()+"/"+f.submitName+".*",
		func(server *server.Server, ctx *server.RequestContext) error {
			parameters := GetParametersFromPath(ctx.VirtualPath)

			for _, parameterMetadata := range f.parametersMetadata {
				if !parameterMetadata.Required {
					continue
				}

				_, ok := parameters[parameterMetadata.Name]
				if !ok {
					return server.SendError(
						ctx.Conn,
						fmt.Sprintf("Missing the '%s' parameter.", parameterMetadata.Name),
					)
				}
			}

			return f.submitCallback(parameters, server, ctx)
		},
	)
}

func (f *Form) Apply(router *server.Router) error {
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

func Apply(
	router *server.Router,
	prefixPath string,
	parameterMetadata []ParameterMetadata,
	submitCallback SubmitCallback,
) error {
	form := NewForm(prefixPath, parameterMetadata, submitCallback)

	err := form.Apply(router)
	if err != nil {
		return err
	}

	return nil
}
