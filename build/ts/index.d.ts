declare interface HTMLElement {
  innerHTML: string;
}

declare interface IFetchResponse {
  status: number;
  json: () => Promise<object>;
}

declare function fetch(request: Request): Promise<IFetchResponse>;
declare function fetch(url: string, params?: object): Promise<IFetchResponse>;
