<template>
  <div id="login" :class="{ recaptcha: recaptcha }">
    <form @submit="submit">
      <img :src="logoURL" alt="File Browser" />
      <h1>{{ name }}</h1>
      <div v-if="error !== ''" class="wrong">{{ error }}</div>

      <input
        autofocus
        class="input input--block"
        type="text"
        autocapitalize="off"
        v-model="username"
        :placeholder="t('login.username')"
      />
      <input
        class="input input--block"
        type="password"
        v-model="password"
        :placeholder="t('login.password')"
      />
      <input
        class="input input--block"
        v-if="createMode"
        type="password"
        v-model="passwordConfirm"
        :placeholder="t('login.passwordConfirm')"
      />

      <div v-if="recaptcha" id="recaptcha"></div>
      <input
        class="button button--block"
        type="submit"
        :value="createMode ? t('login.signup') : t('login.submit')"
      />

      <p @click="toggleMode" v-if="signup">
        {{ createMode ? t("login.loginInstead") : t("login.createAnAccount") }}
      </p>
      <p v-if="oauth2Config.useOauth2" @click="oauth2">
        {{ oauth2Config.Oauth2Name }}
      </p>
    </form>
  </div>
</template>

<script setup lang="ts">
import { StatusError } from "@/api/utils";
import * as auth from "@/utils/auth";
import {
  name,
  logoURL,
  recaptcha,
  recaptchaKey,
  signup,
} from "@/utils/constants";
import { inject, onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { useRoute, useRouter } from "vue-router";
import { getOauth2LoginState } from "@/api/settings.js";

// Define refs
const createMode = ref<boolean>(false);
const error = ref<string>("");
const username = ref<string>("");
const password = ref<string>("");
const passwordConfirm = ref<string>("");
const oauth2Config = ref({
  useOauth2: false,
  Oauth2Name: "",
  Oauth2Url: "",
});

const route = useRoute();
const router = useRouter();
const { t } = useI18n({});
// Define functions
const toggleMode = () => (createMode.value = !createMode.value);

const $showError = inject<IToastError>("$showError")!;

const submit = async (event: Event) => {
  event.preventDefault();
  event.stopPropagation();

  const redirect = (route.query.redirect || "/files/") as string;

  let captcha = "";
  if (recaptcha) {
    captcha = window.grecaptcha.getResponse();

    if (captcha === "") {
      error.value = t("login.wrongCredentials");
      return;
    }
  }

  if (createMode.value) {
    if (password.value !== passwordConfirm.value) {
      error.value = t("login.passwordsDontMatch");
      return;
    }
  }

  try {
    if (createMode.value) {
      await auth.signup(username.value, password.value);
    }

    await auth.login(username.value, password.value, captcha);
    router.push({ path: redirect });
  } catch (e: any) {
    // console.error(e);
    if (e instanceof StatusError) {
      if (e.status === 409) {
        error.value = t("login.usernameTaken");
      } else if (e.status === 403) {
        error.value = t("login.wrongCredentials");
      } else {
        $showError(e);
      }
    }
  }
};

const updateOauthData = async () => {
  const data = await getOauth2LoginState()
  oauth2Config.value.useOauth2 = data["use"];
  if (oauth2Config.value.useOauth2) {
    oauth2Config.value.Oauth2Name = data["name"]
    oauth2Config.value.Oauth2Url = data["url"]
  }
}
updateOauthData()
const oauth2 = () => {
  window.location.href = oauth2Config.value.Oauth2Url;
};

const code = ref();
async function init() {
  if (window.location.search.indexOf("code=") !== -1) {
    // console.log(window.location.search.split("?code=")[1])
    const codes = window.location.search.split("code=");
    code.value = codes[codes.length - 1].split("&")[0];
    let redirect = ""
    try {
      if (window.location.search.indexOf("redirect=") !== -1) {
        const redirects = window.location.search.split("redirect=")
        redirect = redirects[redirects.length-1].split("&")[0].split("?")[0]
      }
    }catch (e) {
      redirect = "/files/";
    }
    if (redirect === "" || redirect === undefined || redirect === null) {
      redirect = "/files/";
    }
    try {
      await auth.oauth2callback(code.value);
      await router.push({ path: redirect });
    } catch (e) {
      if (e.message == 409) {
        $showError(t("login.usernameTaken"));
        router.push({ name: "Login" });
      } else {
        $showError(t("login.wrongCredentials"));
        router.push({ name: "Login" });
      }
    }
  }
}
init();
// Run hooks
onMounted(() => {
  if (!recaptcha) return;

  window.grecaptcha.ready(function () {
    window.grecaptcha.render("recaptcha", {
      sitekey: recaptchaKey,
    });
  });
});

</script>
