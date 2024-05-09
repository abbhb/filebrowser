<template>
  <errors v-if="error" :errorCode="error.status" />
  <div class="row" v-else-if="!layoutStore.loading">
    <div class="column">
      <div class="card">
        <div class="card-title">
          <h2>{{ t("settings.users") }}</h2>
          <input
            type="file"
            ref="fileInput"
            accept=".json"
            style="display: none"
          />

          <button style="margin-left: 3px" class="button" @click="importData">
            <i class="fas fa-file-import"></i>导入
          </button>
          <button style="margin-left: 3px" class="button" @click="exportData">
            <i class="fas fa-file-export"></i>导出
          </button>
          <router-link to="/settings/users/new"
            ><button class="button">
              {{ t("buttons.new") }}
            </button></router-link
          >
        </div>

        <div class="card-content full">
          <table>
            <tr>
              <th>{{ t("settings.username") }}</th>
              <th>{{ t("settings.admin") }}</th>
              <th>{{ t("settings.scope") }}</th>
              <th></th>
            </tr>

            <tr v-for="user in users" :key="user.id">
              <td>{{ user.username }}</td>
              <td>
                <i v-if="user.perm.admin" class="material-icons">done</i
                ><i v-else class="material-icons">close</i>
              </td>
              <td>{{ user.scope }}</td>
              <td class="small">
                <router-link :to="'/settings/users/' + user.id"
                  ><i class="material-icons">mode_edit</i></router-link
                >
              </td>
            </tr>
          </table>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useLayoutStore } from "@/stores/layout";
import { users as api } from "@/api";
import Errors from "@/views/Errors.vue";
import { onMounted, ref } from "vue";
import { useI18n } from "vue-i18n";
import { StatusError } from "@/api/utils";
import { postUserSync } from "@/api/sync.js";

const error = ref<StatusError | null>(null);
const users = ref<IUser[]>([]);
const fileInput = ref();
const layoutStore = useLayoutStore();
const { t } = useI18n();

const uploadFile = async (file: any) => {
  const formData = new FormData();
  formData.append("file", file);
  const res = await postUserSync(formData);

  if (res.status === 200) {
    alert("导入成功，请刷新页面");
    window.location.reload();
  } else {
    throw new Error(String(res.status));
  }
};
const importData = () => {
  // 处理导入数据的逻辑
  console.log("Importing data...");
  fileInput.value.click();
  fileInput.value.addEventListener("change", () => {
    const file = fileInput.value.files[0];
    if (file) {
      if (file.size > 10 * 1024 * 1024) {
        // 检查文件大小，10MB以下
        alert("文件大小不能超过10MB");
        return;
      }
      if (file.type !== "application/json") {
        // 检查文件类型
        alert("请选择JSON文件");
        return;
      }
      uploadFile(file);
    }
  });
};
const exportData = () => {
  // 处理导出数据的逻辑
  window.open("/api/sync/user");
  // console.log("Exporting data...");
  // const data = await getUserSync();
  // const jsond = data.json();
  // const jsonData = JSON.stringify(jsond);
  // // 创建Blob对象
  // const blob = new Blob([jsonData], { type: "application/json" });
  // // 创建URL指向Blob对象
  // const url = window.URL.createObjectURL(blob);
  // // 创建<a>元素
  // const a = document.createElement("a");
  // // 设置<a>元素的href为URL
  // a.href = url;
  // // 设置下载文件的名称
  // a.download = "users.json";
  // // 将<a>元素隐藏在页面中
  // a.style.display = "none";
  // // 将<a>元素添加到DOM中
  // document.body.appendChild(a);
  // // 模拟点击<a>元素以下载文件
  // a.click();
  // // 释放URL对象
  // window.URL.revokeObjectURL(url);
};
onMounted(async () => {
  layoutStore.loading = true;

  try {
    users.value = await api.getAll();
  } catch (err) {
    if (err instanceof Error) {
      error.value = err;
    }
  } finally {
    layoutStore.loading = false;
  }
});
</script>
