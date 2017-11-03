var State = (function () {
    function State() {
    }
    State.prototype.fetchState = function () {
        var _this = this;
        var request = new Request("/state", {
            method: "GET",
            headers: new Headers({
                "Accept": "application/json"
            })
        });
        fetch(request)
            .then(function (response) {
            if (response.status != 200) {
                console.error("Rubbernecker responded with non 200 http status.");
            }
            return response.json();
        })
            .then(function (data) {
            _this.content = data;
            $(document).trigger(State.updated);
        });
    };
    return State;
}());
State.updated = "rubbernecker:state:updated";
var Application = (function () {
    function Application() {
        this.state = new State();
    }
    Application.prototype.dealCard = function (card) {
        var tmpl = document.getElementById("card-template");
        var parsed = document.createElement("div");
        if (tmpl === null) {
            console.error("No card-template provided!");
            return;
        }
        parsed.innerHTML = tmpl.innerHTML;
        var $card = $(parsed).find("> div");
        $card.attr("style", "display: none;");
        this.updateCardData($card, card);
        $("#" + card.status)
            .append(parsed.innerHTML);
        this
            .gracefulIn($("#" + card.status + " #" + card.id));
    };
    Application.prototype.gracefulIn = function ($elements) {
        var _this = this;
        $elements.each(function () {
            var $element = $(_this);
            if (!$element.is(":hidden")) {
                return;
            }
            $element.css("opacity", 0);
            $element.slideDown();
            setTimeout(function () {
                $element.animate({
                    opacity: 1
                });
            }, 750);
        });
    };
    Application.prototype.gracefulOut = function ($elements) {
        var _this = this;
        $elements.each(function () {
            var $element = $(_this);
            if ($element.is(":hidden")) {
                return;
            }
            $element.css("opacity", 1);
            $element.animate({
                opacity: 0
            });
            setTimeout(function () {
                $element.slideUp();
            }, 750);
        });
    };
    Application.prototype.listAssignees = function (card) {
        var assignees = [];
        $.each(Object.keys(card.assignees), function (i, id) {
            assignees.push("<li>" + card.assignees[id].name + "</li>");
        });
        return assignees.join("");
    };
    Application.prototype.parseContent = function () {
        var _this = this;
        if (typeof this.state.content.cards === "undefined") {
            console.error("No cards found in state...");
            return;
        }
        var cards = this.state.content.cards;
        for (var i in cards) {
            var card = cards[i];
            var $card = $("#" + card.id);
            if (typeof $card !== "undefined") {
                this.updateCard($card, card);
            }
            else {
                this.dealCard(card);
            }
        }
        this.updateFreeMembers(this.state.content.free_team_members);
        $.each(Object.keys(this.state.content.support), function (i, type) {
            return $("body > header ." + type).text(_this.state.content.support[type].member);
        });
        $(document).trigger(Application.updated);
    };
    Application.prototype.run = function () {
        var _this = this;
        console.info("Running rubbernecker application.");
        setInterval(function () {
            _this.state.fetchState();
        }, 15000);
        $(document)
            .on(State.updated, function () { _this.parseContent(); });
    };
    Application.prototype.setAssignees = function ($card, card) {
        var html;
        if (Object.keys(card.assignees).length > 0) {
            html = "<h4>Assignee" + (Object.keys(card.assignees).length > 1 ? "s" : "") + "</h4>\n        <ul>" + this.listAssignees(card) + "</ul>";
        }
        else {
            html = "<h4 class=\"text-danger\">Nobody is working on this</h4>\n        <p>Sad times.</p>";
        }
        $card
            .find("> main")
            .html(html);
        return this;
    };
    Application.prototype.setHeader = function ($card, card) {
        $card
            .find("> header > a")
            .attr("href", card.url)
            .text(card.title);
        $card
            .find("> header > span")
            .text(card.in_play + " day" + (card.in_play !== 1 ? "s" : ""));
        return this;
    };
    Application.prototype.setStickers = function ($card, card) {
        return this;
    };
    Application.prototype.setupCard = function ($card, card) {
        $card
            .attr("class", "card " + card.status)
            .attr("id", card.id);
        return this;
    };
    Application.prototype.updateCard = function ($card, card) {
        var correctState = $card.parents("#" + card.status).length > 0;
        if (!correctState) {
            setTimeout(function () {
                $card.remove();
            }, 500);
            this.gracefulOut($card);
            this.dealCard(card);
        }
        else {
            this.updateCardData($card, card);
        }
    };
    Application.prototype.updateCardData = function ($card, card) {
        this.setupCard($card, card)
            .setHeader($card, card)
            .setAssignees($card, card)
            .setStickers($card, card);
    };
    Application.prototype.updateFreeMembers = function (freeMembers) {
        var $freeMembers = $("body > footer");
        $freeMembers
            .find("span")
            .text(Object.keys(freeMembers).length);
        $freeMembers.find("ul").empty();
        $.each(Object.keys(freeMembers), function (i, id) {
            return $freeMembers.find("ul").append("<li>" + freeMembers[id].name + "</li>");
        });
    };
    return Application;
}());
Application.updated = "rubbernecker:application:updated";
