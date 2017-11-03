type Members = { [id: string]: Member }

interface Member {
  id: number
  email: string
  name: string
}

interface Card {
  id: number
  assignees: Members
  in_play: number
  status: string
  stickers: any
  title: string
  url: string
}

interface Support {
  type: string
  member: string
}

interface Response {
  cards: Card[]
  support: { [type: string]: Support }
  free_team_members: Members
}

declare var $: any;

class State {
  content: Response;

  static updated: string = "rubbernecker:state:updated";

  constructor() {
    this.content = new Response();
  }

  fetchState() {
    let request = new Request("/state", {
      method: "GET",
      headers: new Headers({
        "Accept": "application/json",
      })
    });

    fetch(request)
      .then(response => {
        if (response.status != 200) {
          console.error("Rubbernecker responded with non 200 http status.");
        }

        return response.json();
      })
      .then(data => {
        this.content = data;

        // Trigger updated event.
        $(document).trigger(State.updated);
      });
  }
}

class Application {
  static updated: string = "rubbernecker:application:updated";

  private state: State;

  constructor() {
    this.state = new State();
  }

  private dealCard(card: Card) {
    let tmpl: HTMLElement | null = document.getElementById("card-template");
    let parsed = document.createElement("div");

    if (tmpl === null) {
      console.error("No card-template provided!");
      return
    }

    parsed.innerHTML = tmpl.innerHTML;

    let $card = $(parsed).find("> div");

    $card.attr("style", "display: none;");

    this.updateCardData($card, card);

    $("#" + card.status)
      .append(parsed.innerHTML);

    this
      .gracefulIn($("#" + card.status + " #" + card.id))
  }

  gracefulIn($elements: any) {
    $elements.each(() => {
      let $element = $(this);

      if (!$element.is(":hidden")) {
        return;
      }

      $element.css("opacity", 0);
      $element.slideDown();

      setTimeout(() => {
        $element.animate({
          opacity: 1
        });
      }, 750);
    });
  }

  gracefulOut($elements: any) {
    $elements.each(() => {
      let $element = $(this);

      if ($element.is(":hidden")) {
        return;
      }

      $element.css("opacity", 1);
      $element.animate({
        opacity: 0
      });

      setTimeout(() => {
        $element.slideUp();
      }, 750);
    });
  }

  listAssignees(card: Card) {
    let assignees: string[] = [];

    $.each(Object.keys(card.assignees), (i: number, id: number) => {
      assignees.push(`<li>` + card.assignees[id].name + `</li>`);
    });

    return assignees.join("");
  }

  private parseContent() {
    if (typeof this.state.content.cards === "undefined") {
      console.error("No cards found in state...");
      return
    }

    let cards = this.state.content.cards;

    for (let i in cards) {
      let card = cards[i];

      let $card = $("#" + card.id);


      if (typeof $card !== "undefined") {
        this.updateCard($card, card);
      } else {
        this.dealCard(card);
      }
    }

    this.updateFreeMembers(this.state.content.free_team_members);

    $.each(Object.keys(this.state.content.support), (i: number, type: string) =>
      $("body > header ." + type).text(this.state.content.support[type].member)
    );

    $(document).trigger(Application.updated);
  }

  run() {
    console.info("Running rubbernecker application.");

    setInterval(() => {
      this.state.fetchState();
    }, 15000);

    $(document)
      .on(State.updated, () => { this.parseContent() });
  }

  private setAssignees($card: any, card: Card) {
    let html;

    if (Object.keys(card.assignees).length > 0) {
      html = `<h4>Assignee` + (Object.keys(card.assignees).length > 1 ? `s` : ``) + `</h4>
        <ul>` + this.listAssignees(card) + `</ul>`;
    } else {
      html = `<h4 class="text-danger">Nobody is working on this</h4>
        <p>Sad times.</p>`;
    }

    $card
      .find("> main")
      .html(html);

    return this;
  }

  private setHeader($card: any, card: Card) {
    $card
      .find("> header > a")
      .attr("href", card.url)
      .text(card.title);

    $card
      .find("> header > span")
      .text(card.in_play + " day" + (card.in_play !== 1 ? "s" : ""));

    return this;
  }

  private setStickers($card: any, card: Card) {


    return this;
  }

  private setupCard($card: any, card: Card) {
    $card
      .attr("class", "card " + card.status)
      .attr("id", card.id);

    return this;
  }

  private updateCard($card: any, card: Card) {
    let correctState = $card.parents("#" + card.status).length > 0;

    if (!correctState) {
      setTimeout(() => {
        $card.remove();
      }, 500);

      this.gracefulOut($card);
      this.dealCard(card);
    } else {
      this.updateCardData($card, card);
    }
  }

  private updateCardData($card: any, card: Card) {
    this.setupCard($card, card)
      .setHeader($card, card)
      .setAssignees($card, card)
      .setStickers($card, card);
  }

  private updateFreeMembers(freeMembers: Members) {
    let $freeMembers = $("body > footer");

    $freeMembers
      .find("span")
      .text(Object.keys(freeMembers).length);

    $freeMembers.find("ul").empty();

    $.each(Object.keys(freeMembers), (i: number, id: string) =>
      $freeMembers.find("ul").append("<li>" + freeMembers[id].name + "</li>")
    );
  }
}
