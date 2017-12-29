"use strict";

class Eventer {
    constructor() {
        this.listeners = {};
    }

    addEventListener(name, f) {
        if (!this.listeners[name])
            this.listeners[name] = [];
        this.listeners[name].push(f);
    }

    removeEventListener(name, f) {
        if (this.listeners[name]) {
            this.listeners[name] =
                this.listeners[name].filter(v => v != f);
        }
    }

    dispatchEvent(name, o) {
        let toCall = this.listeners[name];
        if (toCall)
            for (let idx in toCall)
                toCall[idx](o);
        return true;
    }
}

function xhrPromise(method, endPoint, body) {
    return new Promise(function(res, rej) {
        let xhr = new XMLHttpRequest();
        xhr.onload = () => { res(xhr.responseText); }
        xhr.onerror = () => { rej(xhr.responseText); }
        xhr.open(method, endPoint);
        xhr.send(body);
    }); 
}

function loadJson(endPoint) {
    return xhrPromise("GET", endPoint, "").then(JSON.parse);
}

class IdNamePairSelector {
    constructor(all, curr) {
        this.element = document.createElement("select");
        for (let id in all) {
            let opt = document.createElement("option");
            opt.value = id;
            opt.innerHTML = all[id];
            opt.selected = id == curr;
            this.element.appendChild(opt);
        }
    }

    asId() {
        let options = this.element.getElementsByTagName("option");
        for (let idx in options)
            if (options[idx].selected)
                return options[idx].value|0;
    }
}

function applyStyle(elem, style) {
    if (style)
        for (let key in style)
            elem.style[key] = style[key];
}

class NumberInput {
    constructor(value) {
        this.element = document.createElement("input");
        this.element.type = "number";
        this.element.step = "any";
        this.element.value = value;
    }

    asFloat() {
        return parseFloat(this.element.value);
    }

    asInt() {
        return this.element.value|0;
    }
}

class Labeled {
    constructor(label, comp, style, spanStyle) {
        this.element = document.createElement("div");
        let span = document.createElement("span");
        span.innerHTML = label;
        applyStyle(span, spanStyle);
        applyStyle(this.element, style);
        this.element.appendChild(span);
        this.element.appendChild(comp.element);
    }
}

class EditorSelector {
    constructor(gameInfo, sources, editorClass, newSource) {
        let container = document.createElement("div");
        let sourceContainer = document.createElement("div");
        let select = document.createElement("select");
        let emptyOption = document.createElement("option");
        emptyOption.innerHTML = "";
        emptyOption.value = -1;
        let newOption = document.createElement("option");
        newOption.innerHTML = "[New]";
        select.appendChild(emptyOption);
        select.appendChild(newOption);
        for (let i in sources) {
            let option = document.createElement("option");
            option.innerHTML = sources[i].Name;
            option.value = i;
            select.appendChild(option);
        }
        select.addEventListener('change', ev => {
            if (ev.target.value == -1)
                return;
            let source = sources[ev.target.value];
            if (!source) {
                let srcId = 0;
                for (let idx in sources)
                    if (sources[idx].Id >= srcId)
                        srcId = sources[idx].Id + 1;
                source = newSource(srcId);
            }
            let display = new editorClass(gameInfo, source);
            display.addEventListener('save', source => {
                sources[source.Id] = source;
                if (ev.target.value > 0) {
                    ev.target.value = source.Id;
                    ev.target.innerHTML = source.Name;
                } else {
                    let option = document.createElement("option");
                    option.innerHTML = source.Name;
                    option.value = source.Id;
                    select.appendChild(option);
                    option.selected = true;
                }
            });
            sourceContainer.innerHTML = "";
            sourceContainer.appendChild(display.element);
        });

        container.appendChild(select);
        container.appendChild(sourceContainer);
        this.element = container;
    }
}

function loadGameInfo() {
    function fillIdNamePairs(dict, pairs) {
        for (let idx = 0; idx < pairs.length; idx++)
            dict[pairs[idx].Id] = pairs[idx].Name;
    }
    let gameInfo = {
        damageTypes: {},
        statuses: {},
        scalingFuncs: {},
        stats: {},
        traits: {},
    };
    return loadJson("/staticGameInfo").then(function (gi) {
        fillIdNamePairs(gameInfo.damageTypes, gi.DamageTypes);
        fillIdNamePairs(gameInfo.statuses, gi.Statuses);
        fillIdNamePairs(gameInfo.scalingFuncs, gi.ScalingFuncs);
        fillIdNamePairs(gameInfo.stats, gi.Stats);
        fillIdNamePairs(gameInfo.traits, gi.Traits);
        return gameInfo; 
    });
}

