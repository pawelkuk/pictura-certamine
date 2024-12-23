import _ from "lodash";
import "./style.css";
import "@shoelace-style/shoelace/dist/themes/light.css";
import "filepond/dist/filepond.min.css";
import "@shoelace-style/shoelace/dist/themes/dark.css";
import SlButton from "@shoelace-style/shoelace/dist/components/button/button.js";
import SlIcon from "@shoelace-style/shoelace/dist/components/icon/icon.js";
import SlInput from "@shoelace-style/shoelace/dist/components/input/input.js";
import SlCheckbox from "@shoelace-style/shoelace/dist/components/checkbox/checkbox.js";
import SlRating from "@shoelace-style/shoelace/dist/components/rating/rating.js";
import { setBasePath } from "@shoelace-style/shoelace/dist/utilities/base-path.js";
import * as FilePond from "filepond";

// Set the base path to the folder you copied Shoelace's assets to
setBasePath("/dist/shoelace");

// Get a reference to the file input element
const inputElement = document.querySelector('input[type="file"]');

// Create a FilePond instance
const pond = FilePond.create(inputElement, {
  allowDrop: true,
  allowBrowse: true,
  allowRemove: true,
  allowMultiple: true,
  id: "art-piece",
  className: "art-piece-file-upload",
  credits: false,
});
