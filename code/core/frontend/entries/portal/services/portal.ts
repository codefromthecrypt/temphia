import type { SiteUtils } from "../../../lib/utils/site";
import { ApiManager } from "./apm";
import { Navigator } from "./nav";
import { Notifier } from "./notifier";

export interface AppOptions {
  base_url: string;
  tenant_id: string;
  user_token: string;
  site_utils: SiteUtils;
}

export interface Utils {
  toast_success(msg: string): void;
  toast_error(msg: string): void;
  big_modal_open(compo: any, options: object): void;
  big_modal_close(): void;
  small_modal_open(compo: any, options: object): void;
  small_modal_close(): void;
  notification_toggle(): void;
}

export class PortalService {
  options: AppOptions;

  nav: Navigator;
  toaster: any;

  api_manager: ApiManager;
  notifier: Notifier;
  utils: Utils;

  constructor(opts: AppOptions) {
    this.options = opts;
    this.nav = new Navigator();
  }

  async init() {
    this.api_manager = new ApiManager(
      this.options.base_url,
      this.options.tenant_id,
      this.options.user_token
    );
    await this.api_manager.init();

    this.notifier = new Notifier(this.api_manager.self_api);
    await this.notifier.init()
  }

  inject(utils: Utils) {
    this.utils = utils;
  }
}