define("types", ["require", "exports"], function (require, exports) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
});
define("state", ["require", "exports", "tslib", "jquery"], function (require, exports, tslib_1, jquery_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    jquery_1 = tslib_1.__importDefault(jquery_1);
    var State = (function () {
        function State() {
            this.content = {
                cards: [],
                support: {},
                free_team_members: {},
            };
        }
        State.prototype.fetchState = function () {
            return tslib_1.__awaiter(this, void 0, void 0, function () {
                var request, response, data;
                return tslib_1.__generator(this, function (_a) {
                    switch (_a.label) {
                        case 0:
                            request = new Request('/state', {
                                method: 'GET',
                                headers: new Headers({
                                    Accept: 'application/json',
                                }),
                            });
                            return [4, fetch(request)];
                        case 1:
                            response = _a.sent();
                            if (response.status !== 200) {
                                console.error('Rubbernecker responded with non 200 http status.');
                            }
                            return [4, response.json()];
                        case 2:
                            data = _a.sent();
                            this.content = data;
                            jquery_1.default(document).trigger(State.updated);
                            return [2];
                    }
                });
            });
        };
        State.updated = 'rubbernecker:state:updated';
        return State;
    }());
    exports.default = State;
});
define("app", ["require", "exports", "tslib", "jquery", "state"], function (require, exports, tslib_2, jquery_2, state_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    jquery_2 = tslib_2.__importDefault(jquery_2);
    state_1 = tslib_2.__importDefault(state_1);
    var Application = (function () {
        function Application() {
            this.state = new state_1.default();
        }
        Application.prototype.dealCard = function (card) {
            var tmpl = document.getElementById('card-template');
            var parsed = document.createElement('div');
            if (tmpl === null) {
                console.error('No card-template provided!');
                return;
            }
            parsed.innerHTML = tmpl.innerHTML;
            var $card = jquery_2.default(parsed).find('> div');
            $card.attr('style', 'display: none;');
            this.updateCardData($card, card);
            jquery_2.default("#" + card.status)
                .append(parsed.innerHTML);
            this
                .gracefulIn(jquery_2.default("#" + card.status + " #" + card.id));
        };
        Application.prototype.gracefulIn = function ($elements) {
            var _this = this;
            $elements.each(function () {
                var $element = jquery_2.default(_this);
                if (!$element.is(':hidden')) {
                    return;
                }
                $element.css('opacity', 0);
                $element.slideDown();
                setTimeout(function () {
                    $element.animate({
                        opacity: 1,
                    });
                }, 750);
            });
        };
        Application.prototype.gracefulOut = function ($elements) {
            var _this = this;
            $elements.each(function () {
                var $element = jquery_2.default(_this);
                if ($element.is(':hidden')) {
                    return;
                }
                $element.css('opacity', 1);
                $element.animate({
                    opacity: 0,
                });
                setTimeout(function () {
                    $element.slideUp();
                }, 750);
            });
        };
        Application.prototype.listAssignees = function (card) {
            var assignees = [];
            jquery_2.default.each(Object.keys(card.assignees), function (_, id) {
                assignees.push("<li>" + card.assignees[id].name + "</li>");
            });
            return assignees.join('');
        };
        Application.prototype.run = function () {
            var _this = this;
            console.info('Running rubbernecker application.');
            setInterval(function () {
                _this.state.fetchState();
            }, 15000);
            jquery_2.default(document)
                .on(state_1.default.updated, function () { _this.parseContent(); });
        };
        Application.prototype.parseContent = function () {
            var _this = this;
            if (typeof this.state.content.cards === 'undefined') {
                console.error('No cards found in state...');
                return;
            }
            var cards = this.state.content.cards;
            for (var _i = 0, cards_1 = cards; _i < cards_1.length; _i++) {
                var card = cards_1[_i];
                var $card = jquery_2.default("#" + card.id);
                if ($card) {
                    this.updateCard($card, card);
                }
                else {
                    this.dealCard(card);
                }
            }
            this.updateFreeMembers(this.state.content.free_team_members);
            jquery_2.default.each(Object.keys(this.state.content.support), function (_, schedule) {
                return jquery_2.default("body > header ." + schedule).text(_this.state.content.support[schedule].member);
            });
            jquery_2.default(document).trigger(Application.updated);
        };
        Application.prototype.setAssignees = function ($card, card) {
            var html;
            if (Object.keys(card.assignees).length > 0) {
                html = "<h4>Assignee" + (Object.keys(card.assignees).length > 1 ? "s" : "") + "</h4>\n        <ul>" + this.listAssignees(card) + "</ul>";
            }
            else {
                html = "<h4 class='text-danger'>Nobody is working on this</h4>\n        <p>Sad times.</p>";
            }
            $card
                .find('> main')
                .html(html);
            return this;
        };
        Application.prototype.setHeader = function ($card, card) {
            $card
                .find('> header > a')
                .attr('href', card.url)
                .text(card.title);
            $card
                .find('> header > span')
                .text(card.in_play + " day" + (card.in_play !== 1 ? 's' : ''));
            return this;
        };
        Application.prototype.setStickers = function ($card, card) {
            console.log($card, card);
            return this;
        };
        Application.prototype.setupCard = function ($card, card) {
            $card
                .attr('class', "card " + card.status)
                .attr('id', card.id);
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
            var $freeMembers = jquery_2.default('body > footer');
            $freeMembers
                .find('span')
                .text(Object.keys(freeMembers).length);
            $freeMembers.find('ul').empty();
            jquery_2.default.each(Object.keys(freeMembers), function (_, id) {
                return $freeMembers.find('ul').append("<li>" + freeMembers[id].name + "</li>");
            });
        };
        Application.updated = 'rubbernecker:application:updated';
        return Application;
    }());
    exports.Application = Application;
});
define("index", ["require", "exports", "app"], function (require, exports, app_1) {
    "use strict";
    Object.defineProperty(exports, "__esModule", { value: true });
    var app = new app_1.Application();
    app.run();
});
