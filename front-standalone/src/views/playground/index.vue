<script setup lang="ts">
import { computed, onMounted, reactive, ref } from "vue";
import Authorization from './components/Authorization.vue'
import Password from './components/Password.vue'
import Client from './components/Client.vue'
import { fetchConfig } from '/@/api/common'
import { Setting } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { LocalStorageService } from "/@/utils/persistence"

defineOptions({
  name: "Playground"
});

const grantTypes = ref("1");
const configData = reactive({
  trust_domain: [],
  authorization_endpoint: "",
  token_endpoint: "",
  userinfo_endpoint: "",
  default_scope: "",
  access_token_type: "bearer",
  client_id: "",
  client_secret: ""
});

const att = [
  {label: "authorization header w/ Bearer prefix", value: "bearer"},
  {label: "access_token URL parameter", value: "query"}
];

async function getGlobalConfig() {
  const {data} = await fetchConfig();
  Object.assign(configData, data);
  const lss = new LocalStorageService();
  const id = lss.getItem("id");
  const secret = lss.getItem("secret");
  Object.assign(configData, {client_id: id??"", client_secret: secret??""});
  lss.removeItem("id");
  lss.removeItem("secret");
}

function handleSaveTokenType(){
  ElMessage.success(configData.access_token_type);
}

function handleSaveAccount() {
  ElMessage.success(configData.client_id);
}

const computedPopWidth = computed(() => {
  return window.innerWidth < 640 ? "100%" : "60%";
});

onMounted(() => {
  getGlobalConfig();
})
</script>

<template>
  <el-main>
    <div class="mainContent">
      <div class="rowHeader">
        <el-row>
          <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8" class="titleRow">
            <span class="bigTitle">OAuth 2.0 Playground</span>
          </el-col>
          <el-col :xs="20" :sm="22" :md="14" :lg="14" :xl="14" class="contentRow">
            <el-radio-group v-model="grantTypes" size="small">
              <el-radio-button label="1">Authorization Code</el-radio-button>
              <el-radio-button label="2">Resource Owner Password Credentials</el-radio-button>
              <el-radio-button label="3">Client Credentials</el-radio-button>
            </el-radio-group>
          </el-col>
          <el-col :xs="4" :sm="2" :md="2" :lg="2" :xl="2" class="contentRow">
            <div id="settings" style="position: fixed;right: 0">
              <el-popover placement="bottom" trigger="click" :width="computedPopWidth">
                <template #reference>
                  <el-button style="vertical-align: middle;">
                    <el-icon>
                      <Setting/>
                    </el-icon>
                  </el-button>
                </template>
                <el-descriptions title="Config" :column="1" border>
                  <el-descriptions-item label="Trust Domain">
                    <el-tag v-for="item in configData.trust_domain" :key="item" style="margin-right: 5px">
                      {{ item }}
                    </el-tag>
                  </el-descriptions-item>
                  <el-descriptions-item label="Authorization Endpoint">
                    {{ configData.authorization_endpoint }}
                  </el-descriptions-item>
                  <el-descriptions-item label="Token Endpoint">
                    {{ configData.token_endpoint }}
                  </el-descriptions-item>
                  <el-descriptions-item label="Userinfo Endpoint">
                    {{ configData.userinfo_endpoint }}
                  </el-descriptions-item>
                  <el-descriptions-item label="Default Scope">
                    {{ configData.default_scope }}
                  </el-descriptions-item>
                  <el-descriptions-item label="Client Id">
                    <el-input size="small" v-model="configData.client_id" placeholder="Client Id"/>
                  </el-descriptions-item>
                  <el-descriptions-item label="Client Secret">
                    <el-input size="small" v-model="configData.client_secret" placeholder="Client Secret"/>
                  </el-descriptions-item>
                  <el-descriptions-item label="Access token location">
                    <el-select
                        v-model="configData.access_token_type"
                        size="small"
                        style="width: 100%"
                        @change="handleSaveTokenType"
                    >
                      <el-option
                          v-for="item in att"
                          :key="item.label"
                          :label="item.label"
                          :value="item.value"
                      />
                    </el-select>
                  </el-descriptions-item>
                </el-descriptions>
              </el-popover>
            </div>
          </el-col>
        </el-row>
      </div>
      <div class="rowContent">
        <el-row v-show="grantTypes==='1'" class="contentBox">
          <Authorization :cfgData="configData"></Authorization>
        </el-row>
        <el-row v-show="grantTypes==='2'" class="contentBox">
          <Password :cfgData="configData"></Password>
        </el-row>
        <el-row v-show="grantTypes==='3'" class="contentBox">
          <Client :cfgData="configData"></Client>
        </el-row>
      </div>
    </div>
  </el-main>
</template>
<style lang="less" scoped>

:deep(.el-radio-button) {
  --el-radio-button-checked-bg-color: #b70031;
  --el-radio-button-checked-text-color: var(--el-color-white);
  --el-radio-button-checked-border-color: #b70031;
  --el-radio-button-disabled-checked-fill: var(--el-border-color-extra-light);

  .el-radio-button__inner:hover {
    color: #b70031;
    background-color: #fff;
  }

  .el-radio-button__inner:active {
    color: #fff;
    background-color: #b70031;
  }

  .el-radio-button__inner:focus {
    color: #fff;
    background-color: #b70031;
  }

  .el-radio-button__original-radio:checked + .el-radio-button__inner {
    color: #fff;
    background-color: #b70031;
    border-color: #b70031;
    box-shadow: #b70031;
  }
}

.el-main {
  padding: 0;
  height: calc(100vh - 60px);

  .mainContent {
    height: 100%;

    .rowHeader {
      width: 100%;
      height: 60px;

      .titleRow {
        height: 60px;
        line-height: 60px;
        text-align: left;

        .bigTitle {
          color: #b70031;
          padding: 0 20px;
          font-size: 20px;
          font-weight: bolder;
        }
      }

      .contentRow {
        text-align: left;
        height: 60px;
        line-height: 60px;
      }

      .modeSelect {
      }
    }

    .rowContent {
      width: 100%;
      height: calc(100% - 60px);

      .contentBox {
        height: 100%;
      }
    }
  }
}

@media only screen and (max-width: 768px) {
  .el-main {
    width: 100%;

    .mainContent {
      height: 100%;

      .rowHeader {
        width: 100%;
        height: 80px;

        .titleRow {
          height: 40px;
          line-height: 40px;

          .bigTitle {
            color: #dd4b39;
            padding: 0 20px;
            font-size: 18px;
            font-weight: 300;
          }
        }

        .contentRow {
          height: 40px;
          line-height: 30px;
          padding: 0 20px;
          text-align: left;
        }
      }
    }
  }
}
</style>

