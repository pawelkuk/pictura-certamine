import _ from "lodash";
import "./style.css";
import "@shoelace-style/shoelace/dist/themes/light.css";
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
