// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package view

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func BlueBox() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"blue-box\"><h1>CUM SĂ PARTICIPAȚI?</h1><p>Achiziționați un produs (e) Marvel de la firma Hasbro, in Carrefour. Creați cu orice tehnică o lucrare interesantă pentru concurs, care să descrie eroul dvs preferat de la Marvel (poate fi un desen, un videoclip, o înregistrare audio, o descriere etc.). Înregistrați-o, făcând clic pe butonul Participați de mai jos. Apoi va trebui să completați datele dvs și să adăugați lucrarea și, la  final, să citiți și să confirmați acceptarea <a href=\"/assets/img/REGULILE_CONCURSULUI.pdf\">întregului regulament</a>. </p><span onclick=\"location.href=&#39;/?dialog=open&#39;;\" class=\"button\"></span><p class=\"small\">Competiția  <strong>este deschisă consumatorilor finali cu domiciliul și/sau reședința în România, care achiziționează\u00a0un produs Marvel de la firma Hasbro, din magazinele Carrefour. Lista completă a produselor care condiționează participarea la concurs și metoda de achiziție pot fi găsite în Regulament.</strong> Pentru a participa, vi se va cere să furnizați  datele  dvs. și să adăugați lucrarea de concurs și,  la final, să citiți și să confirmați acceptarea <a href=\"/assets/img/REGULILE_CONCURSULUI.pdf\">întregului regulament</a>. În cazul în care lucrarea de concurs este depusă de un minor, acesta trebuie să aibă consimțământul reprezentantului legal sau al tutorelui legal pentru a participa la concurs, pe care, în cazul unei victorii, va fi obligat să îl trimită în conformitate cu Regulamentul (conținutul consimțământului necesar este Anexa 1 la Regulament).</p><p class=\"small\">Concursul va începe pe 27 ianuarie 2025 la ora 00:00 și se va încheia pe 6 martie 2025 la ora 23:59. Comisia de concurs va selecta câștigătorul până la 31 martie 2025. Înainte de a revendica premiul, câștigătorul va trebui să furnizeze detaliile dovezii achiziției (rețineți că trebuie să păstrați dovada originală a achiziției). Vom contacta câștigătorul prin e-mail. Nu uitați să verificați întotdeauna secțiunea Spam a căsuței poștale electronice. </p><br><p class=\"small\">Valoarea premiului este de 7.433,38 euro, la care se adaugă TVA, în conformitate cu prevederile legislației române.**</p><br><p class=\"small\">Aici puteți citi <a href=\"/assets/img/REGULILE_CONCURSULUI.pdf\" target=\"_blank\">regulamentul complet</a> și <a href=\"/assets/img/politica_de_confidentialitate.pdf\" target=\"_blank\">politica de confidențialitate</a> a concursului. </p></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
