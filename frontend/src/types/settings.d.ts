import { str } from "video.js";

export interface ISettings {
  signup: boolean;
  createUserDir: boolean;
  userHomeBasePath: string;
  defaults: SettingsDefaults;
  rules: any[];
  branding: SettingsBranding;
  tus: SettingsTus;
  shell: string[];
  commands: SettingsCommand;
  oauth2:Oauth2
}

export interface SettingsDefaults {
  scope: string;
  locale: string;
  viewMode: ViewModeType;
  singleClick: boolean;
  sorting: Sorting;
  perm: Permissions;
  commands: any[];
  hideDotfiles: boolean;
  dateFormat: boolean;
}
export interface Oauth2 {
  disable?:boolean;
  name?:string;
  tokenurl?:string;
  meurl?:string;
  userinfourl?:string;
  clientid?:string;
  clientsecret?:string;
  redirecturi?:string;
  scope?:string;
  state?:string;
  authorizeurl?:string;
}

export interface SettingsBranding {
  name: string;
  disableExternal: boolean;
  disableUsedPercentage: boolean;
  files: string;
  theme: UserTheme;
  color: string;
}

export interface SettingsTus {
  chunkSize: number;
  retryCount: number;
}

export interface SettingsCommand {
  after_copy?: string[];
  after_delete?: string[];
  after_rename?: string[];
  after_save?: string[];
  after_upload?: string[];
  before_copy?: string[];
  before_delete?: string[];
  before_rename?: string[];
  before_save?: string[];
  before_upload?: string[];
}

export interface SettingsUnit {
  KB: number;
  MB: number;
  GB: number;
  TB: number;
}
