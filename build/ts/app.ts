import State from './state';
import { ICard, IMembers } from './types';

declare var $: any;

export class Application {
  public static updated: string = 'rubbernecker:application:updated';

  private state: State;

  constructor() {
    this.state = new State();
  }

  public dealCard(card: ICard) {
    const tmpl: HTMLElement | null =
      document.getElementById(card.status === 'done' ? 'thin-card-template' : 'card-template');
    const parsed = document.createElement('div');

    if (tmpl === null) {
      console.error('No card-template provided!');
      return;
    }

    parsed.innerHTML = tmpl.innerHTML; // tslint:disable-line:no-inner-html

    const $card = $(parsed).find('> div');

    $card.attr('style', 'display: none;');

    this.updateCardData($card, card);

    $(`#${card.status}`)
      .append(parsed.innerHTML);

    this
      .gracefulIn($(`#${card.status} #${card.id}`));
  }

  public gracefulIn($elements: any) {
    $elements.each(() => {
      const $element = $(this);

      if (!$element.is(':hidden')) {
        return;
      }

      $element.css('opacity', 0);
      $element.slideDown();

      setTimeout(() => {
        $element.animate({
          opacity: 1,
        });
      }, 750);
    });
  }

  public gracefulOut($elements: any) {
    $elements.each(() => {
      const $element = $(this);

      if ($element.is(':hidden')) {
        return;
      }

      $element.css('opacity', 1);
      $element.animate({
        opacity: 0,
      });

      setTimeout(() => {
        $element.slideUp();
      }, 750);
    });
  }

  public run() {
    console.info('Running rubbernecker application.');

    setInterval(() => {
      this.state.fetchState();
    }, 15000);

    $(document)
      .on(State.updated, () => { this.parseContent(); });
  }

  private parseContent() {
    if (typeof this.state.content.cards === 'undefined') {
      console.error('No cards found in state...');
      return;
    }

    const cards = this.state.content.cards;

    for (const card of cards) {
      const $card = $(`#${card.id}`);

      if ($card) {
        this.updateCard($card, card);
      } else {
        this.dealCard(card);
      }
    }

    this.updateFreeMembers(this.state.content.free_team_members);

    $.each(Object.keys(this.state.content.support), (_: number, schedule: string) =>
      $(`body > header .${schedule}`).text(this.state.content.support[schedule].member),
    );

    $(document).trigger(Application.updated);
  }

  private setAssignees($card: any, card: ICard) {
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

  private setHeader($card: any, card: ICard) {
    $card
      .find('> h3 > a')
      .attr('href', card.url)
      .text(card.title);

    $card
      .find('> footer > .elapsed > small')
      .text(`${card.in_play} day${card.in_play !== 1 ? 's' : ''}`);

    return this;
  }

  private setStickers($card: any, card: ICard) {
    const $stickers = $card
      .find('footer > .stickers');

    if ($stickers.length > 0) {
      $stickers
        .empty();

      for (const sticker of card.stickers) {
        $stickers
          .append(`<div><img src="${sticker.Image}" alt="${sticker.Title}" title="${sticker.Title}"></div>`);
      }
    }

    return this;
  }

  private setupCard($card: any, card: ICard) {
    $card
      .attr('class', `card ${card.status}`)
      .attr('id', card.id);

    return this;
  }

  private updateCard($card: any, card: ICard) {
    const correctState = $card.parents(`#${card.status}`).length > 0;

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

  private updateCardData($card: any, card: ICard) {
    this.setupCard($card, card)
      .setHeader($card, card)
      .setAssignees($card, card)
      .setStickers($card, card);
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
