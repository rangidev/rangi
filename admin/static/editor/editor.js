// Text
class RangiText extends HTMLParagraphElement {
    constructor() {
        self = super();
    }

    connectedCallback() {
        this.contentEditable = "true";
    }
}
customElements.define("rangi-text", RangiText, { extends: "p" });

// Title
class RangiTitle extends HTMLHeadingElement {
    constructor() {
        self = super();
    }

    connectedCallback() {
        this.contentEditable = "true";
        this.classList.add("h1");
    }
}
customElements.define("rangi-title", RangiTitle, { extends: "h1" });

// Reference
class RangiReference extends HTMLButtonElement {
    constructor() {
        self = super();
    }

    connectedCallback() {
        this.innerHTML = "Reference Button";
        this.classList.add("btn", "btn-lg", "btn-primary");
    }
}
customElements.define("rangi-reference", RangiReference, { extends: "button" });