// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.747
package view

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func PrizeSection() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"prize\"><h1>W konkursie możesz wygrać:</h1><img src=\"/assets/img/prize.png\"><h1>Nagroda będzie można zrealizować do 1.02.2026 roku. Wartość nagrody to 40000 zł**.</h1><div class=\"conditions\"><p>*Warunki korzystania z nagrody: </p><ul><li>Z pobytu można skorzystać do 1 lutego 2026r. </li><li>Daty wyłączone z korzystania z nagrody to: </li><li>Podróż należy zarezerwować co najmniej 30 dni przed datą podróży. </li><li>Daty i ustalenia dotyczące podróży zależą od dostępności, ograniczeń podróży lotniczych, świąt oraz innych ograniczeń dotyczących nagród i podróży. </li><li>Wśród użytkowników wycieczki przynajmniej jedna osoba musi mieć ukończone 18 lat. </li><li>Nagroda nie obejmuje przejazdu z domu na lotnisko i odwrotnie, noclegu w pobliżu lotniska przed odlotem i po przylocie, ani niczego, co nie zostało wyraźnie określone. </li><li>Wszelkie wydatki inne niż wskazane tutaj będą wyłączną odpowiedzialnością zwycięzcy i jego/jej osób towarzyszących. </li><li>Wyłączone są w szczególności wycieczki, ubezpieczenie podróżne, usługi pralnicze, zabiegi wellness, specjalne napoje alkoholowe i bezalkoholowe, indywidualne przekąski, rozmowy telefoniczne krajowe i międzynarodowe, łącze internetowe, gadżety i nieobowiązkowe napiwki. </li><li>Wszyscy użytkownicy podróży muszą podróżować w tych samych terminach, tymi samymi lotami i przebywać w tym samym pokoju w hotelu. </li><li>Loty w klasie ekonomicznej będą odbywać się w zależności od dostępności lotu w datach podanych w momencie potwierdzenia. Wszelkie podatki podróżne, których nie można przewidzieć, zostaną opłacone na miejscu przez zwycięzcę. </li><li>Doba dotyczy czterech osób (przy przeliczeniu na 4 osoby dorosłe) w pokoju czteroosobowym. </li><li>Hotel wymaga preautoryzacji karty kredytowej lub depozytu w gotówce na pokrycie kosztów osobistych (obsługa pokoju, opłaty dodatkowe, minibar itp.). Wszelkie podatki lokalne należy uiścić bezpośrednio w hotelu na koszt zwycięzcy. </li><li>Zwycięzca i jego osoba towarzysząca muszą posiadać i mieć przy sobie, na swoją wyłączną opiekę i koszt, ważne dwowody osobiste lub paszporty oraz wszelkie wizy/zezwolenia/zezwolenia [„niezależnie od nazwy, niezbędne do podróży”]. </li><li>Zwycięzca i jego osoby towarzyszące muszą przez cały czas przestrzegać odpowiednich ograniczeń i przepisów dotyczących zdrowia i bezpieczeństwa związanych z wirusem Covid-19, mających zastosowanie do podróży, hoteli i zajęć. Mamy tu na myśli, przykładowo i nie wyłącznie, testy i/lub inne środki potwierdzające status szczepienia, jeśli takie będą wymagane. </li><li>Aby zakwalifikować się do nagrody , za każdorazowe przestrzeganie takich i podobnych ograniczeń i/lub wymagań wyłączną odpowiedzialność ponosi zwycięzca i jego osoba towarzysząca . Z wyjątkiem przypadków przewidzianych dla trzech osób towarzyszących objętych nagrodą, zwycięzcy nie wolno zabierać ze sobą na wycieczkę z nagrodą dodatkowych krewnych lub osób towarzyszących. Zwycięzca i jego towarzysze muszą podróżować tą samą trasą z jednego z lotnisk wybranych przez organizatora. Po wybraniu towarzyszy nie będzie możliwa żadna zamiana bez wyraźnej zgody promotora, z zastrzeżeniem wyłącznie jego uznania. </li><li>Zastrzegamy sobie prawo do zastąpienia wydarzenia podobną usługą w przypadku niedostępności. </li><li>Zastrzegamy sobie prawo do zmiany szczegółów nagrody w przypadku okoliczności od nas niezależnych, w tym między innymi: zmiany harmonogramu lotu, dostępności hotelu, wojny, zamieszek, zamieszek i niepokojów społecznych, nieodpowiednich warunków klimatycznych, zwiększonej siły, dostępności atrakcji, zezwoleń rządowych lub pandemie. </li><li>W przypadku sytuacji wskazanej powyżej zwycięzca będzie musiał zaproponować alternatywne dwie możliwe daty wyjazdu, z których najbliższa data będzie terminem oddzielonym o co najmniej 30 dni od daty przesłania propozycji dat alternatywnych. W wybrane daty nie można uwzględnić świąt państwowych. </li><li>Dołożymy wszelkich starań, aby wyjazd odbył się w żądanym terminie orientacyjnym, jednak nie jesteśmy w stanie zagwarantować terminów letnich. Święta i daty zbiegające się z konkretnymi wydarzeniami nie są dostępne. </li><li>Nagroda jest nieprzenoszalna, nie podlega zwrotowi i nie może zostać zamieniona na gotówkę; zwycięzca nie może zmienić, zmienić, zastąpić ani przedłużyć żadnego elementu nagrody (w całości lub w części). </li><li>Określono, że zwycięzca nie otrzyma żadnego zwrotu pieniędzy, jeśli w momencie rezerwacji wartość nagrody będzie niższa niż szacunkowa wartość wskazana powyżej. </li><li>Zwycięzca otrzyma gwarancję usług wskazanych powyżej i niczego więcej nie może oczekiwać. </li></ul><p>**Wartość wskazanych nagród należy rozumieć według stanu na dzień powstania regulaminu. </p></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return templ_7745c5c3_Err
	})
}
