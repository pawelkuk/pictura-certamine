import _ from "lodash";
import "@shoelace-style/shoelace/dist/themes/light.css";
import "vanilla-cookieconsent/dist/cookieconsent.css";
import "filepond/dist/filepond.min.css";
import "@shoelace-style/shoelace/dist/themes/dark.css";
import SlButton from "@shoelace-style/shoelace/dist/components/button/button.js";
import SlIcon from "@shoelace-style/shoelace/dist/components/icon/icon.js";
import SlIconButton from "@shoelace-style/shoelace/dist/components/icon-button/icon-button.js";
import SlInput from "@shoelace-style/shoelace/dist/components/input/input.js";
import SlAlert from "@shoelace-style/shoelace/dist/components/alert/alert.js";
import SlCheckbox from "@shoelace-style/shoelace/dist/components/checkbox/checkbox.js";
import SlRating from "@shoelace-style/shoelace/dist/components/rating/rating.js";
import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
import * as FilePond from "filepond";
import FilePondPluginFileValidateType from "filepond-plugin-file-validate-type";
import FilePondPluginFileValidateSize from "filepond-plugin-file-validate-size";
import * as CookieConsent from "vanilla-cookieconsent";
import { DataTable } from "simple-datatables";
import "simple-datatables/dist/style.css";
import "./style.css";
import "../img/background-sky.png";
import "../img/banner.png";
import "../img/banner_mobile.png";
import "../img/disney_logo.png";
import "../img/disney_paris.png";
import "../img/disneyland_day.png";
import "../img/iron_man_night.png";
import "../img/marvel.png";
import "../img/marvel_apla.png";
import "../img/marvel_left.png";
import "../img/marvel_ramka.png";
import "../img/marvel_right.png";
import "../img/prize.png";
import "../img/qrcode.png";
import "../img/disney-background-web.webp";
import "../img/hasbro.svg";
import "../img/informatii_privind_prelucrarea_datelor_cu_caracter_personal.pdf";
import "../img/politica_de_confidentialitate.pdf";
import "../img/REGULILE_CONCURSULUI.pdf";
// Set the base path to the folder you copied Shoelace's assets to
setBasePath("/assets/shoelace");
FilePond.registerPlugin(FilePondPluginFileValidateType);
FilePond.registerPlugin(FilePondPluginFileValidateSize);
// Get a reference to the file input element
const inputElement = document.querySelector('input[type="file"]');

// Create a FilePond instance
FilePond.create(inputElement, {
  allowFileSizeValidation: true,
  maxFileSize: "128MB",
  maxTotalFileSize: "128MB",
  acceptedFileTypes: [
    "image/png",
    "image/jpeg",
    "image/jpg",
    "image/gif",
    "audio/wav",
    "audio/mpeg",
    "image/tiff",
    "application/pdf",
    "video/mp4",
    "application/msword",
    "audio/x-wma",
    "audio/x-ms-wma",
    "audio/x-m4a",
    ".wma",
  ],
  fileValidateTypeDetectType: (source, type) =>
    new Promise((resolve, reject) => {
      if (/\.wma$/.test(source.name)) return resolve("audio/x-ms-wma");

      // accept detected type
      resolve(type);
    }),
  allowDrop: true,
  allowBrowse: true,
  allowRemove: true,
  allowMultiple: true,
  storeAsFile: true,
  id: "art-piece",
  name: "art-piece",
  className: "art-piece-file-upload",
  multiple: true,
  credits: false,
  // server: "/upload",
  required: true,
});

CookieConsent.run({
  categories: {
    necessary: {
      enabled: true, // this category is enabled by default
      readOnly: true, // this category cannot be disabled
    },
    analytics: {},
    performance: {},
    functional: {},
  },
  // onChange: ({ cookie }) => {
  //   console.log("on change");
  // },
  onConsent: ({ cookie }) => {
    const cookieBannerBlur = document.getElementById("cookie-banner-blur");
    if (cookieBannerBlur != null) {
      cookieBannerBlur.style.display = "none";
    }
  },
  // onFirstConsent: ({ cookie }) => {
  //   console.log("on first consent");
  // },
  // onModalHide: ({ cookie }) => {
  //   console.log("on modal hide");
  // },
  // onModalShow: ({ cookie }) => {
  //   console.log("on modal show");
  // },
  onModalReady: ({ cookie }) => {
    const cookieBannerBlur = document.getElementById("cookie-banner-blur");
    if (cookieBannerBlur != null) {
      cookieBannerBlur.style.display = "block";
    }
  },
  language: {
    default: "ro",
    translations: {
      ro: {
        consentModal: {
          title: "Informații despre cookie-uri",
          description: `Împreună cu partenerii noștri, procesăm informații despre dvs., dispozitivele dvs. și comportamentul dvs. online,  utilizând tehnologii precum cookie-urile, pentru a furniza, analiza și îmbunătăți serviciile noastre, pentru a personaliza conținutul sau în scopuri publicitare pe acest site și pe alte site-uri web, pe aplicații sau platforme și pentru a furniza servicii de social media. Pentru a afla mai multe, vă rugăm să consultați Politica noastră <a href="/assets/img/politica_de_confidentialitate.pdf" target="_blank">privind cookie-urile</a>.`,
          acceptAllBtn: "Acceptați totul",
          acceptNecessaryBtn: "Renunțați",
          showPreferencesBtn: "Managementul preferințelor",
        },
        preferencesModal: {
          title: "Centru de preferințe cookie",
          acceptAllBtn: "Acceptați totul",
          acceptNecessaryBtn: "Renunțați",
          savePreferencesBtn: "Acceptați selecția curentă",
          closeIconLabel: "Close modal",
          sections: [
            {
              title:
                "Gestionarea preferințelor de consimțământ pentru prelucrarea datelor",
              description:
                "Deoarece vă respectăm dreptul la confidențialitate, puteți alege să nu vă dați acordul pentru anumite tipuri de cookie-uri. Faceți clic pe fiecare titlu de categorie pentru a afla mai multe și pentru a modifica setările implicite. Vă rugăm să rețineți că blocarea anumitor tipuri de cookie-uri poate afecta negativ utilizarea site-ului nostru web și a serviciilor pe care le putem oferi. Utilizarea acestui instrument va duce la setarea unui cookie pe dispozitivul dvs. pentru a memora preferințele pe care le puteți modifica în orice moment. Pentru a afla mai multe, vă rugăm să consultați Politica noastră privind cookie-urile.",
            },
            {
              title: "Cookie-urile strict necesare sunt întotdeauna active",
              description:
                "Aceste cookie-uri sunt esențiale pentru funcționarea site-ului web și pentru funcționalitățile sale de bază și nu pot fi dezactivate în sistemele noastre. De obicei, acestea sunt setate doar ca răspuns la acțiunile dvs., care constituie o solicitare de acțiune, cum ar fi setarea preferințelor de confidențialitate, accesarea, căutarea sau descoperirea conținutului, completarea formularelor sau trimiterea de conținut. Puteți seta browserul să vă blocheze sau să vă avertizeze despre aceste cookie-uri, dar in acest caz, unele părți ale site-ului nu vor funcționa.",
              //this field will generate a toggle linked to the 'necessary' category
              linkedCategory: "necessary",
            },
            {
              title: "Cookie-uri de performanță",
              description:
                "Aceste cookie-uri ne permit să numărăm vizitele și sursele de trafic, astfel încât să putem măsura și îmbunătăți performanța site-ului nostru (inclusiv Google Analytics). Acestea ne ajută să știm care pagini sunt cele mai populare și cum navighează vizitatorii pe site. Puteți permite aceste cookie-uri și vă puteți retrage permisiunea în orice moment.",
              linkedCategory: "performance",
            },
            {
              title: "Cookie-uri funcționale",
              description:
                "Folosite de noi pentru a detecta sau a memora alegerile pe care le faceți pentru a vă personaliza utilizarea site-ului nostru web, cum ar fi limba, locația sau alte setări. Puteți permite aceste cookie-uri și vă puteți retrage permisiunea în orice moment. Dezactivarea acestor cookie-uri poate afecta modul în care funcționează site-ul web.",
              linkedCategory: "functional",
            },
            {
              title: "Cookie-uri de profilare și publicitate",
              description:
                "Aceste cookie-uri pot fi setate prin intermediul site-ului nostru de către noi și/sau partenerii noștri de publicitate. Acestea pot fi utilizate pentru a crea un profil al preferințelor dumneavoastră și pentru a afișa reclame relevante pe acest site și pe alte site-uri. Nu pot să stocheze în mod direct informații personale, ci se pot baza pe identificarea unică a browserului și a dispozitivului dvs. de internet. Puteți permite aceste cookie-uri și vă puteți retrage permisiunea în orice moment. Dezactivarea acestor cookie-uri nu afectează funcționarea site-ului web.",
              linkedCategory: "analytics",
            },
          ],
        },
      },
    },
  },
});

const t = document.querySelector("#crm-contest-list");
if (t != null) {
  new DataTable(t, {
    searchable: true,
    fixedHeight: false,
    perPage: 100,
    perPageSelect: [10, 100, 1000],
  });
}

const showContestForm = document.getElementById("participate");
const contestFormDialog = document.getElementById("contest-form-dialog");
if (showContestForm != null) {
  const url = new URL(window.location.href);
  if (url.searchParams.get("dialog") != null) {
    contestFormDialog.showModal();
  }
}
