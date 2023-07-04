<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { generateRandomString } from "/@/utils/x";
import { LocalStorageService } from "/@/utils/persistence"
import useClipboard from 'vue-clipboard3';
import { fetchACToken, fetchApiData, fetchRefreshToken } from "/@/api/playground";

const props = defineProps({
  cfgData: {
    type: Object,
    default: () => {
      return {
        trust_domain: [],
        authorization_endpoint: "",
        token_endpoint: "",
        userinfo_endpoint: "",
        default_scope: "",
        access_token_type: "",
        client_id: "",
        client_secret: ""
      };
    }
  },
});
const {toClipboard} = useClipboard();
const isHorizontal = ref(false);
const toggleLayout = () => {
  isHorizontal.value = !isHorizontal.value;
};

// Http Info
const requestInfo = reactive({
  meta: {
    method: '',
    path: '',
    proto: '',
    host: '',
  },
  header: '',
  body: '',
  code: '',
});
const responseInfo = reactive({
  meta: {
    method: '',
    path: '',
    proto: '',
    host: '',
    status: '',
  },
  header: '',
  strHeader: '',
  body: '',
  code: '',
});
const rawJsonInfo = reactive({});
const exampleInfo = reactive({});
const isWrapReq = ref(true);//控制body是否自动换行
const isWrapRes = ref(true);//控制body是否自动换行

function updateReqAndRes() {
  requestInfo.body = atob(requestInfo.body);
  requestInfo.code = JSON.stringify(requestInfo.header, undefined, 0.1)
  //去掉首尾的{}
  requestInfo.code = requestInfo.code.slice(1, requestInfo.code.length - 1);
  requestInfo.code = requestInfo.meta.method + ' ' + requestInfo.meta.path + ' ' + requestInfo.meta.proto + requestInfo.code;
  requestInfo.code = requestInfo.code + 'Host: ' + requestInfo.meta.host;
  requestInfo.code = requestInfo.code.replaceAll('"', '');

  // responseInfo.meta = res["data"]["response"]["meta"];
  // responseInfo.body = res["data"]["response"]["body"];
  // responseInfo.body = atob(responseInfo.body);
  responseInfo.body = decodeURIComponent(escape(window.atob(responseInfo.body)))
  // responseInfo.header = res["data"]["response"]["header"];
  responseInfo.strHeader = JSON.stringify(responseInfo.header, undefined, 0.1)
  //去掉首尾的{}
  responseInfo.strHeader = responseInfo.strHeader.slice(1, responseInfo.strHeader.length - 1)
  responseInfo.strHeader = responseInfo.meta.proto + ' ' + responseInfo.meta.status + responseInfo.strHeader;
  responseInfo.strHeader = responseInfo.strHeader.replaceAll('"', '')
}

// Step 1
const activeName = ref('1');
const s1Data = reactive({
  authorization_endpoint: "",
  redirect_uri: window.location.href.split("?")[0],
  scope: "",
  response_type: "code",
  state: "",
});

const initialAddress = ref("");

function handleS1Change() {
  initialAddress.value = s1Data.authorization_endpoint.concat(
      "?response_type=code",
      s1Data.scope?.length > 0 ? "&scope=".concat(s1Data.scope) : "",
      props.cfgData.client_id?.length > 0 ? "&client_id=".concat(props.cfgData.client_id) : "",
      "&redirect_uri=",
      s1Data.redirect_uri,
      s1Data.state?.length > 0 ? "&state=".concat(s1Data.state) : ""
  );
}

function handleGetAuthorizationCode() {
  if(s1Data.scope.length === 0){
    ElMessage.error('scope不能为空');
    return;
  }else if(s1Data.state.length === 0){
    ElMessage.error('state不能为空');
    return;
  }else{
    if(props.cfgData.client_id.length === 0){
      ElMessage.error('client_id is empty, please click the config button on the right side, and check the configuration');
      return;
    }else{
      const lss = new LocalStorageService();
      const ci = {key: "id", value: props.cfgData.client_id};
      const cs = {key: "secret", value: props.cfgData.client_secret};
      lss.addItem(ci);
      if(cs.value.length > 0){
        lss.addItem(cs);
      }
      window.location.href = initialAddress.value;

    }
  }
}

// Step 2
const code = ref("");
const state = ref("");
const currentToken = ref("");
const currentRefreshToken = ref("");

async function handleGetToken() {
  if(code.value.length === 0){
    ElMessage.error('code不能为空');
    return;
  }else if(props.cfgData.client_id.length === 0){
    ElMessage.error('client_id不能为空');
    return;
  }else if(props.cfgData.client_secret.length === 0){
    ElMessage.error('client_secret不能为空');
    return;
  }else{
    const dataObject = {
      code: code.value,
      client_id: props.cfgData.client_id,
      client_secret: props.cfgData.client_secret,
      scope: props.cfgData.default_scope,
      redirect_uri: window.location.href.split("?")[0]
    };
    fetchACToken(dataObject).then(({code, msg, data}) => {
      if(code === 0){
        const {request, response, rawjson, example} = data;
        const {access_token, refresh_token} = rawjson;
        currentToken.value = access_token;
        currentRefreshToken.value = refresh_token;
        s3CurrentToken.value = access_token;
        toClipboard(access_token).finally(() => {
          ElMessage.success(`已复制access_token: ${access_token}`);
        });
        Object.assign(requestInfo, request);
        Object.assign(responseInfo, response);
        Object.assign(rawJsonInfo, rawjson);
        Object.assign(exampleInfo, example);

        window.history.replaceState({}, document.title, window.location.pathname);
        updateReqAndRes();
      }else{
        ElMessage.error(msg);
      }
    });

  }
}

function handleRefreshToken() {
  if(props.cfgData.client_id.length === 0){
    ElMessage.error('client_id is empty, please click the config button on the right side, and check the configuration');
    return;
  }else if(props.cfgData.client_secret.length === 0){
    ElMessage.error('client_secret is empty, please click the config button on the right side, and check the configuration');
    return;
  }else if(currentRefreshToken.value.length === 0){
    ElMessage.error('refresh_token is empty, please get the access_token firstly');
    return;
  }else{
    const dataObject = {
      refresh_token: currentRefreshToken.value,
      client_id: props.cfgData.client_id,
      client_secret: props.cfgData.client_secret
    };

    fetchRefreshToken(dataObject).then(({code, msg, data}) => {
      if(code === 0){
        const {request, response, rawjson, example} = data;
        const {access_token, refresh_token} = rawjson;
        currentToken.value = access_token;
        currentRefreshToken.value = refresh_token;
        s3CurrentToken.value = access_token;
        toClipboard(access_token).finally(() => {
          ElMessage.success(`已复制新access_token: ${access_token}`);
        });
        Object.assign(requestInfo, request);
        Object.assign(responseInfo, response);
        Object.assign(rawJsonInfo, rawjson);
        Object.assign(exampleInfo, example);
        updateReqAndRes();
      }else{
        ElMessage.error(msg);
        return;
      }
    });
  }
}

// Step 3
const requestUri = ref("");
const requestMethod = ref("GET");
const s3CurrentToken = ref("");
const s3TokenType = ref("Bearer");

const methodOptions = [
  {label: 'GET', value: 'GET', disabled: false},
  {label: 'POST', value: 'POST', disabled: true},
  {label: 'PUT', value: 'PUT', disabled: true},
  {label: 'DELETE', value: 'DELETE', disabled: true},
  {label: 'PATCH', value: 'PATCH', disabled: true},
  {label: 'HEAD', value: 'HEAD', disabled: true},
  {label: 'OPTIONS', value: 'OPTIONS', disabled: true}
];

function handleMethodChange(value) {
  requestMethod.value = value;
  ElMessage.info(`当前请求方法: ${value}`);
}

async function handleRequestAPI() {
  if(requestUri.value.length === 0){
    ElMessage.error('请求地址不能为空');
    return;
  }else if(s3CurrentToken.value.length === 0){
    ElMessage.error('access_token不能为空');
    return;
  }else{
    const dataObject = {
      method: requestMethod.value,
      api_addr: requestUri.value,
      access_token: s3CurrentToken.value,
      access_token_type: s3TokenType.value,
      header: {},
      http_body: ""
    };
    fetchApiData(dataObject).then(({code, msg, data}) => {
      if(code === 0){
        const {request, response, rawjson, example} = data;
        Object.assign(requestInfo, request);
        Object.assign(responseInfo, response);
        Object.assign(rawJsonInfo, rawjson);
        Object.assign(exampleInfo, example);
        updateReqAndRes();
      }else{
        ElMessage.error(msg);
        return;
      }
    });
  }
}

watch(props.cfgData, (newValue) => {
  s1Data.authorization_endpoint = newValue.authorization_endpoint;
  s1Data.scope = newValue.default_scope;
  initialAddress.value = newValue.authorization_endpoint.concat(
      "?response_type=code",
      newValue.default_scope?.length > 0 ? "&scope=".concat(newValue.default_scope) : "",
      newValue.client_id?.length > 0 ? "&client_id=".concat(newValue.client_id) : "",
      "&redirect_uri=",
      s1Data.redirect_uri,
      s1Data.state?.length > 0 ? "&state=".concat(s1Data.state) : ""
  );
  requestUri.value = newValue.userinfo_endpoint;
  s3TokenType.value = newValue.access_token_type;
});

onMounted(() => {
  // handleDrag(agS1Ref.value.$el, agS1ContainerRef.value.$el);
  s1Data.state = generateRandomString(props.cfgData.default_scope);

  const urlParams = new URLSearchParams(window.location.search);
  if(urlParams.get('code')?.length > 0){
    code.value = urlParams.get('code');
    state.value = urlParams.get('state');
    activeName.value = '2';
    // window.history.replaceState({}, document.title, window.location.pathname);
  }
  requestUri.value = props.cfgData.userinfo_endpoint;
  s3TokenType.value = props.cfgData.access_token_type;

});


const agS1ContainerRef = ref(null);
const agS1Ref = ref(null);
const handleDrag = (floatButton, container) => {
  // todo: drag float button
};

</script>
<template>
  <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
    <div class="demo-collapse" id="main">
      <el-collapse v-model="activeName" accordion>
        <el-collapse-item class="el-collapse-item" name="1">
          <template #title>
            <span class="stepTitle">Step 1: Request for Authorization Code</span>
          </template>
          <el-scrollbar class="fitSide" ref="agS1ContainerRef">
            <h4 style="text-align: left;margin: 0">Authorization Endpoint</h4>
            <el-input v-model="s1Data.authorization_endpoint" disabled/>
            <h4 style="text-align: left;margin: 0">Redirect Uri</h4>
            <el-input v-model="s1Data.redirect_uri" disabled/>
            <h4 style="text-align: left;margin: 0">Scope</h4>
            <el-input v-model="s1Data.scope" placeholder="Scope" @input="handleS1Change" @blur="handleS1Change"/>
            <h4 style="text-align: left;margin: 0">Response Type</h4>
            <el-input v-model="s1Data.response_type" placeholder="Response Type" disabled/>
            <h4 style="text-align: left;margin: 0">State</h4>
            <el-input v-model="s1Data.state" placeholder="State" @input="handleS1Change" @blur="handleS1Change"/>
            <h4 style="text-align: left;margin: 0">Grant Url</h4>
            <el-input v-model="initialAddress" type="textarea" :autosize="{ minRows: 4, maxRows: 6 }"
                      placeholder="Request grant code address"/>
            <el-button ref="agS1Ref" type="primary" @click="handleGetAuthorizationCode" class="n-button side-button">
              GO
            </el-button>
          </el-scrollbar>
        </el-collapse-item>
        <el-collapse-item class="el-collapse-item" name="2">
          <template #title>
            <span class="stepTitle">Step 2: Fetch token with Authorization Code</span>
          </template>
          <el-scrollbar class="fitSide">
            <h4 style="text-align: left;margin: 0">Authorization Code</h4>
            <el-input v-model="code" placeholder="Authorization Code">
              <template #append>
                <el-button type="primary" @click="handleGetToken" class="t-button">
                  Get the token
                </el-button>
              </template>
            </el-input>
            <h4 style="text-align: left;margin: 0">Refresh Token</h4>
            <el-input v-model="currentRefreshToken" placeholder="Refresh Token">
              <template #append>
                <el-button type="primary" @click="handleRefreshToken" class="t-button">
                  Refresh Token
                </el-button>
              </template>
            </el-input>
            <h4 style="text-align: left;margin: 0">Current Token</h4>
            <div class="tokenArea">
              <strong><code>{{ currentToken }}</code></strong>
            </div>
          </el-scrollbar>
        </el-collapse-item>
        <el-collapse-item class="el-collapse-item" name="3">
          <template #title>
            <span class="stepTitle">Step 3: Request to API with the <code
                style="color:#cd3221">access_token</code></span>
          </template>
          <el-scrollbar class="fitSide">
            <h4 style="text-align: left;margin: 0">Request URI</h4>
            <el-input v-model="requestUri" placeholder="Authorization Code"/>
            <h4 style="text-align: left;margin: 0">Method(Currently, only get method is supported.)</h4>
            <el-select v-model="requestMethod" @change="handleMethodChange">
              <el-option v-for="item in methodOptions" :key="item.label" :label="item.label" :value="item.value"
                         :disabled="item.disabled"/>
            </el-select>
            <h4 style="text-align: left;margin: 0">Token</h4>
            <el-input v-model="s3CurrentToken" placeholder="access_token"/>
            <el-button type="primary" @click="handleRequestAPI" class="t-button">
              Fetch Data
            </el-button>
          </el-scrollbar>
        </el-collapse-item>
      </el-collapse>
    </div>
  </el-col>
  <el-col :xs="24" :sm="24" :md="16" :lg="16" :xl="16">
    <div class="http-container" :class="{ 'http-horizontal': isHorizontal }">
      <div class="http-left">
        <el-divider content-position="left" direction="horizontal" class="http-info-type">
          <span>REQUEST INFO</span>
        </el-divider>
        <div class="http-content" style="text-align: start;padding:0px;position: relative;">

          <el-scrollbar class="http-render">
            <highlightjs autodetect :code="requestInfo.code"/>
            <highlightjs :class="{ 'bodyWrap': isWrapRes }" autodetect :code="requestInfo.body"/>
          </el-scrollbar>
          <el-checkbox style="position: absolute;bottom: 30px;left: 20px;" v-model="isWrapReq" label="Wrap Lines"
                       size="large"/>
        </div>
      </div>
      <div class="http-right">
        <el-divider content-position="left" direction="horizontal" class="http-info-type">
          <span>RESPONSE INFO</span>
        </el-divider>
        <div class="http-content" style="text-align: start;padding: 0em;position: relative;">
          <el-scrollbar class="http-render">
            <highlightjs autodetect :code="responseInfo.strHeader"/>
            <highlightjs :class="{ 'bodyWrap': isWrapRes }" autodetect :code="responseInfo.body"/>
          </el-scrollbar>
          <el-checkbox style="position: absolute;bottom: 30px;left: 20px;" v-model="isWrapRes" label="Wrap Lines"
                       size="large"/>
        </div>
      </div>
      <el-button class="http-button http-side-button" @click="toggleLayout">Switch</el-button>
    </div>
  </el-col>
</template>

<style scoped lang="less">
.bodyWrap {
  white-space: pre-line;
  word-break: break-all;
}

:deep(.el-button) {
  margin-top: 10px;
  color: #fff;
  background-color: #b70031;
  border: none;
  border-radius: 0;

  &:focus {
    color: #fff;
    background-color: #b70031;
  }

  &:active {
    color: #fff;
    background-color: #b70031;
  }

  &:hover {
    color: #fff;
    background-color: #e1003c;
  }
}

:deep(.el-collapse-item__content) {
  padding-bottom: 2px;
}

:deep(.el-select) {
  width: 100%;
}

:deep(.el-divider__text) {
  color: #b70031;
  font-weight: 800;
}

#main {
  height: 100%;
  border-right: 1px solid #e5e5e5;

  .el-collapse {
    --el-collapse-border-color: #e5e5e5;
    --el-collapse-header-height: 48px;
    --el-collapse-header-bg-color: #f1f1f1;
    --el-collapse-header-text-color: const(--el-text-color-primary);
    --el-collapse-header-font-size: 13px;
    --el-collapse-content-bg-color: #fefefe;
    --el-collapse-content-font-size: 13px;
    --el-collapse-content-text-color: const(--el-text-color-primary);
    border-top: #e5e5e5;
    border-bottom: #e5e5e5;
  }

  .el-collapse-item {
    height: 100%;
  }

  .stepTitle {
    font-weight: bold;
    padding: 0 20px;
  }

  .fitSide {
    height: calc(100vh - 273px);
    padding: 2px 5px 0;
    position: relative;

    .t-button {
      color: #fff;
      background-color: rgba(183, 0, 49, 1);
      width: 120px;
    }

    .t-button:hover {
      color: #fff;
      background-color: rgba(213, 22, 79, 1);
    }

    .n-button {
      background-color: rgba(183, 0, 49, 0.14);
    }

    .n-button:hover {
      background-color: rgba(183, 0, 49, 1);
    }

    .side-button {
      position: absolute;
      top: 50%;
      transform: translateY(-50%);
      right: 0;
      width: 50px;
      height: 50px;
      border-radius: 50%;
    }

    .tokenArea {
      display: flex;
      align-items: center;
      justify-content: center;
      border: 2px dashed rgb(239, 239, 239);
      width: 99%;
      min-height: 40px;
      height: 40px;
      padding: 10px 0;
    }

  }

}

.http-container {
  position: relative;
  display: flex;
  flex-direction: column;
  height: calc(100vh - 180px);

  .http-button {
    background-color: rgba(183, 0, 49, 0.14);
  }

  .http-button:hover {
    background-color: rgba(183, 0, 49, 1);
  }

  .http-side-button {
    position: absolute;
    top: 50%;
    transform: translateY(-50%);
    right: 0;
    width: 50px;
    height: 50px;
    border-radius: 50%;
  }
}

.http-left,
.http-right {
  height: 100%;
  flex: 1;
  padding: 2px;
}

.http-horizontal .http-left,
.http-horizontal .http-right {
  flex: 1;
}

.http-horizontal {
  flex-direction: row;
}

.http-content {
  height: 100%;
  border: 1px solid #e5e5e5;

  .http-info-type {
  }

  .http-render {
    height: min-content;
  }
}

@media (max-width: 767px) {
  .http-container {
    flex-direction: column;
  }
}

body .el-scrollbar__wrap {
  overflow-x: hidden;
}</style>
