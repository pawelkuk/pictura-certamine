// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package view

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func Navbar() templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<nav class=\"desktop\"><ul class=\"navbar\"><li><img class=\"marvel\" src=\"/assets/img/marvel.png\"></li><li><img class=\"hasbro\" src=\"/assets/img/hasbro.svg\"></li><li id=\"participate\"><a href=\"/?dialog=open\">Participă</a></li><li><a href=\"/assets/img/REGULILE_CONCURSULUI.pdf\" target=\"_blank\">Regulile concursului</a></li><li><a href=\"/assets/img/politica_de_confidentialitate.pdf\" target=\"_blank\">Politica de confidențialitate</a></li></ul></nav><div class=\"mobile\"><img class=\"marvel\" src=\"/assets/img/marvel.png\"> <img class=\"hasbro\" src=\"/assets/img/hasbro.svg\"> <input class=\"menu-btn\" type=\"checkbox\" id=\"menu-btn\" name=\"menu-btn\"> <label class=\"menu-icon\" for=\"menu-btn\"><span class=\"navicon\" aria-label=\"Hamburger menu &#39;icon&#39;\"></span></label><nav class=\"menu\"><a id=\"participate-mobile\" class=\"nav-item\" href=\"/?dialog=open\">Participă</a> <a class=\"nav-item\" href=\"/assets/img/REGULILE_CONCURSULUI.pdf\" target=\"_blank\">Regulile concursului</a> <a class=\"nav-item\" href=\"/assets/img/politica_de_confidentialitate.pdf\" target=\"_blank\">Politica de confidențialitate</a></nav></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
