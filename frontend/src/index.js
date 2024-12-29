import _ from "lodash";
import "@shoelace-style/shoelace/dist/themes/light.css";
import "vanilla-cookieconsent/dist/cookieconsent.css";
import "filepond/dist/filepond.min.css";
import "@shoelace-style/shoelace/dist/themes/dark.css";
import SlButton from "@shoelace-style/shoelace/dist/components/button/button.js";
import SlIcon from "@shoelace-style/shoelace/dist/components/icon/icon.js";
import SlInput from "@shoelace-style/shoelace/dist/components/input/input.js";
import SlAlert from "@shoelace-style/shoelace/dist/components/alert/alert.js";
import SlCheckbox from "@shoelace-style/shoelace/dist/components/checkbox/checkbox.js";
import SlRating from "@shoelace-style/shoelace/dist/components/rating/rating.js";
import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
import * as FilePond from "filepond";
import * as CookieConsent from "vanilla-cookieconsent";
import { DataTable } from "simple-datatables";
import "simple-datatables/dist/style.css";
import "./style.css";
// Set the base path to the folder you copied Shoelace's assets to
setBasePath("/assets/dist/shoelace");

// Get a reference to the file input element
const inputElement = document.querySelector('input[type="file"]');

// Create a FilePond instance
FilePond.create(inputElement, {
  allowDrop: true,
  allowBrowse: true,
  allowRemove: true,
  allowMultiple: true,
  storeAsFile: true,
  id: "art-piece",
  name: "art-piece",
  className: "art-piece-file-upload",
  multiple: true,
  // credits: false,
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
  },

  language: {
    default: "en",
    translations: {
      en: {
        consentModal: {
          title: "We use cookies",
          description: "Cookie modal description",
          acceptAllBtn: "Accept all",
          acceptNecessaryBtn: "Reject all",
          showPreferencesBtn: "Manage Individual preferences",
        },
        preferencesModal: {
          title: "Manage cookie preferences",
          acceptAllBtn: "Accept all",
          acceptNecessaryBtn: "Reject all",
          savePreferencesBtn: "Accept current selection",
          closeIconLabel: "Close modal",
          sections: [
            {
              title: "Somebody said ... cookies?",
              description: "I want one!",
            },
            {
              title: "Strictly Necessary cookies",
              description:
                "These cookies are essential for the proper functioning of the website and cannot be disabled.",

              //this field will generate a toggle linked to the 'necessary' category
              linkedCategory: "necessary",
            },
            {
              title: "Performance and Analytics",
              description:
                "These cookies collect information about how you use our website. All of the data is anonymized and cannot be used to identify you.",
              linkedCategory: "analytics",
            },
            {
              title: "More information",
              description:
                'For any queries in relation to my policy on cookies and your choices, please <a href="#contact-page">contact us</a>',
            },
          ],
        },
      },
    },
  },
});

let dataTable = new DataTable("#crm-contest-list", {
  searchable: true,
  fixedHeight: false,
  perPage: 100,
  perPageSelect: [10, 100, 1000],
});
