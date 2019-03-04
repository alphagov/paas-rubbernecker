interface ICard {
  readonly id: number;
  readonly assignees: IMembers;
  readonly in_play: number;
  readonly status: string;
  readonly stickers: ReadonlyArray<ISticker>;
  readonly title: string;
  readonly url: string;
}

interface IMembers {
  readonly [id: string]: IMember;
}

interface IMember {
  readonly id: number;
  readonly email: string;
  readonly name: string;
}

interface ISupport {
  readonly type: string; //tslint:disable-line:no-reserved-keywords
  readonly member: string;
}

interface ISticker {
  readonly Name: string;
  readonly Label: boolean;
  readonly Regex: string;
  readonly Title: string;
  readonly Image: string;
  readonly Content: string;
  readonly Aliases: ReadonlyArray<string>;
  readonly Class: string;
}

interface IResponse {
  readonly cards: ReadonlyArray<ICard>;
  readonly support: {
    readonly [type: string]: ISupport; // tslint:disable-line:no-reserved-keywords
  };
  readonly free_team_members: IMembers;
}

interface IStateContent {
  readonly cards: ReadonlyArray<ICard>;
  readonly support: {readonly [schedule: string]: ISupport};
  readonly free_team_members: IMembers;
}

class State {
  public static updated: string = 'rubbernecker:state:updated';

  public content: IStateContent;

  constructor() {
    this.content = {
      cards: [],
      support: {},
      free_team_members: {},
    };
  }

  public async fetchState() {
    const request = new Request('/state', {
      method: 'GET',
      headers: new Headers({
        Accept: 'application/json',
      }),
    });

    const response = await fetch(request);

    if (response.status !== 200) {
      console.error('Rubbernecker responded with non 200 http status.');
    }

    this.content = await response.json() as IStateContent;

    // Trigger updated event.
    $(document).trigger(State.updated);
  }
}

const RESET_FILTERS_TIMEOUT_MS = 30 * 60 * 1000;

class Application {
  public static updated: string = 'rubbernecker:application:updated';

  public filterResetTimeout: number;
  private state: State;

  constructor() {
    this.state = new State();
    this.filterResetTimeout = 0;
  }

  public dealCard(card: ICard, pos: number) {
    const tmpl: HTMLElement | null =
      document.getElementById(card.status === 'done' || card.status === 'next' ?
        'thin-card-template' : 'card-template');
    const parsed = document.createElement('div');

    if (tmpl === null) {
      console.error('No card-template provided!');
      return;
    }

    parsed.innerHTML = tmpl.innerHTML; // tslint:disable-line:no-inner-html

    const $card = $(parsed).find('> div');

    $card.attr('style', 'display: none;');

    this.updateCardData($card, card);

    const cardInPos = $(`#${card.status} div.card:eq(${pos})`);
    if (cardInPos.length === 0) {
      $(`#${card.status}`)
        .append(parsed.innerHTML);
    } else {
      cardInPos.before(parsed.innerHTML);
    }

    this
      .gracefulIn($(`#${card.status} #${card.id}`));
  }

  public gracefulIn($elements: JQuery<HTMLElement>) {
    $.each($elements, (_: number, element: HTMLElement) => {
      if (!$(element).is(':hidden')) {
        return;
      }

      $(element).css('opacity', 0);
      $(element).slideDown();

      setTimeout(() => {
        $(element).animate({
          opacity: 1,
        });
      }, 750);
    });
  }

  public gracefulOut($elements: JQuery<HTMLElement>) {
    $.each($elements, (_: number, element: HTMLElement) => {
      if ($(element).is(':hidden')) {
        return;
      }

      $(element).css('opacity', 1);
      $(element).animate({
        opacity: 0,
      });

      setTimeout(() => {
        $(element).slideUp();
      }, 750);
    });
  }

  public run() {
    console.info('Running rubbernecker application.');

    setInterval(async () => {
      await this.state.fetchState();
    }, 15000);

    $(document)
      .on(State.updated, () => { this.parseContent(); });
  }

  public filterTeam(name: string) {
    const anyTeamCards = $('.card:not(:has(.sticker-team))');
    this.gracefulIn(anyTeamCards);

    const teamCards = $('.card:has(.sticker-team)');
    const visibleTeamCards = teamCards.filter(`:has(.sticker-team.team-${name})`);
    const hiddenTeamCards = teamCards.filter(`:not(:has(.sticker-team.team-${name}))`);

    this.gracefulIn(visibleTeamCards);
    this.gracefulOut(hiddenTeamCards);

    this.filterResetTimeout = setTimeout(() => {
      $('input[name=all]').parents('label').trigger('click');
    }, RESET_FILTERS_TIMEOUT_MS);
  }

  public resetFilter() {
    this.gracefulIn($('.card'));
    clearTimeout(this.filterResetTimeout);
  }

  private async parseContent() {
    if (!this.state.content.cards) {
      console.error('No cards found in state...');
      return;
    }

    const cards = this.state.content.cards;

    for (const card of cards) {
      const $card = $(`#${card.id}`);

      let posInColumn = 0;
      for (const c of cards) {
        if (c.id === card.id) {
          break;
        }
        if (c.status === card.status) {
          posInColumn++;
        }
      }

      if ($card) {
        this.updateCard($card, card, posInColumn);
      } else {
        this.dealCard(card, posInColumn);
      }
    }

    setInterval(() => {
      this.updateCounters();
    }, 150);
    this.updateFreeMembers(this.state.content.free_team_members);

    $.each(Object.keys(this.state.content.support), (_: number, schedule: string) =>
      $(`body > header .${schedule}`).text(this.state.content.support[schedule].member),
    );

    $(document).trigger(Application.updated);
  }

  private setAssignees($card: JQuery<HTMLElement>, card: ICard) {
    const $assignees = $card
      .find('> ul');

    if ($assignees.length > 0) {
      $assignees
        .empty();

      for (const assignee of Object.keys(card.assignees)) {
        $assignees
          .append(`<li>${card.assignees[assignee].name}</li>`);
      }
    }

    return this;
  }

  private setHeader($card: JQuery<HTMLElement>, card: ICard) {
    $card
      .find('> h3 > a')
      .attr('href', card.url)
      .text(card.title);

    $card
      .find('> footer > .elapsed > small')
      .text(`${card.in_play} day${card.in_play !== 1 ? 's' : ''}`);

    return this;
  }

  private setStickers($card: JQuery<HTMLElement>, card: ICard) {
    const $stickers = $card.find('footer > .stickers');
    const $labels = $card.find('footer > .labels');

    if ($stickers.length > 0) {
      $stickers.empty();
      $labels.find('.sticker').remove();

      for (const sticker of card.stickers) {
        const stickerClass = sticker.Class !== '' ? ` ${sticker.Class}` : '';
        const classAttribute = `sticker sticker-${sticker.Name}${stickerClass}`;
        let stickerContent = sticker.Image === '' ?
              sticker.Title :
              `<img src="${sticker.Image}" alt="${sticker.Title}" title="${sticker.Title}" />`;

        if (sticker.Content !== '') {
          stickerContent = `${stickerContent} <small>${sticker.Content}</small>`;
        }

        if (!sticker.Label) {
          $stickers
            .append(`<div class="${classAttribute}">${stickerContent}</div> `);
        } else {
          $labels
            .append(`<div class="${classAttribute}">${stickerContent}</div> `);
        }
      }
    }

    return this;
  }

  private setupCard($card: JQuery<HTMLElement>, card: ICard) {
    $card
      .attr('class', `card ${card.status}`)
      .attr('id', card.id);

    return this;
  }

  private updateCard($card: JQuery<HTMLElement>, card: ICard, pos: number) {
    const correctState = $card.parents(`#${card.status}`).length > 0;

    if (!correctState) {
      setTimeout(() => {
        $card.remove();
      }, 500);

      this.gracefulOut($card);
      this.dealCard(card, pos);
    } else {
      this.updateCardData($card, card);
    }
  }

  private updateCardData($card: JQuery<HTMLElement>, card: ICard) {
    this.setupCard($card, card)
      .setHeader($card, card)
      .setAssignees($card, card)
      .setStickers($card, card);
  }

  private async updateCounters() {
    const $sections = $('[data-cards]');

    $.each($sections, (_: number, section: HTMLElement) => {
      const count = $(section).find('div.card').length;
      const limit = parseInt($(section).find('h2 > small').attr('data-limit') || '0', 10);

      $(section).find('h2 > small').removeClass('text-danger');

      $(section).find('h2 > small').attr('data-cards', count);
      $(section).find('h2 > small > span').text(count);

      if (limit !== 0 && count > limit) {
        $(section).find('h2 > small').addClass('text-danger');
      }
    });
  }

  private updateFreeMembers(freeMembers: IMembers) {
    const $freeMembers = $('body > footer');

    $freeMembers
      .find('span')
      .text(Object.keys(freeMembers).length);

    $freeMembers.find('ul').empty();

    $.each(Object.keys(freeMembers), (_: number, id: string) =>
      $freeMembers.find('ul').append(`<li>${freeMembers[id].name}</li>`),
    );
  }
}
