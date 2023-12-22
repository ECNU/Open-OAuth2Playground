<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import {ElMessage, FormInstance} from 'element-plus'
import { LocalStorageService } from "/@/utils/persistence"
import useClipboard from 'vue-clipboard3';
import {fetchACTokenByDevice, fetchApiData, fetchRefreshToken, fetchUserCode} from "/@/api/playground";
import * as QRCode from 'qrcode'
import axios from "axios";

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
const isWrapRes = ref(true);

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
  // redirect_uri: window.location.href.split("?")[0],
  scope: "",
  response_type: "device_code",
  // state: "",
});

const initialAddress = ref("");

// 修改的同时拼接成url显示在Grant Url中
function handleS1Change() {
  initialAddress.value = s1Data.authorization_endpoint.concat(
      "?response_type=device_code",
      s1Data.scope?.length > 0 ? "&scope=".concat(s1Data.scope) : "",
      props.cfgData.client_id?.length > 0 ? "&client_id=".concat(props.cfgData.client_id) : ""
  );
}

function handleDeviceFlow() {
  if (s1Data.scope.length === 0) {
    ElMessage.error('scope cannot be empty');
    return;
  } else {
    if (props.cfgData.client_id.length === 0) {
      ElMessage.error('client_id is empty, please click the config button on the right side, and check the configuration');
      return;
    } else {
      const lss = new LocalStorageService();
      const ci = {key: "id", value: props.cfgData.client_id};
      const cs = {key: "secret", value: props.cfgData.client_secret};
      lss.addItem(ci);
      if(cs.value.length > 0){
        lss.addItem(cs);
      }
      // window.location.href = initialAddress.value;
      const dataObject = {
        initialAddress: initialAddress.value
      };
      fetchUserCode(dataObject).then(({code, msg, data}) => {
        if (code === 0) {
          const {request, response, rawjson, example} = data;
          if (rawjson.user_code === undefined || rawjson.user_code === '')
            return
          // const {interval, verification_uri, user_code, expires_in, device_code} = rawjson || {};
          if (user_code !== undefined && user_code !== null) {
            user_code.value = rawjson.user_code
            verification_uri.value = rawjson.verification_uri
            device_code.value = rawjson.device_code
            expires_in.value = rawjson.expires_in
            // 检查是否获取到user_code,若长度不为0，则跳到Step 2
            if (device_code.value.length === 0 || user_code.value.length === 0 || verification_uri.value.length === 0) {
              ElMessage.error("Authorization failed. Check your configurations.");
              return;
            }
            activeName.value = '2';
            // 轮询是否获取到token
            tokenAvailablelong(expires_in.value);
          }
        } else {
          ElMessage.error("Get user_code failed. Please check your configuration!")
        }
      }).catch((err) => {
        console.error('get user code failed: ', err)
      })
    }
  }
}

// Step 2
const currentToken = ref("");
const currentRefreshToken = ref("");
const verification_uri = ref("");
const user_code = ref("");
const device_code = ref("");
const expires_in = ref("");
const interval = ref("");

const timerId = ref(null);
const qrCodeSrc = ref("");

// 长轮询
function tokenAvailablelong(expire) {
  const dataObject = {
    client_id: props.cfgData.client_id,
    code: device_code.value,
    response_type: "device_code",
    expires_in: expire
  };
  fetchACTokenByDevice(dataObject).then(({code, msg, data}) => {
    if (code === 0) {
      const {request, response, rawjson, example} = data;
      const {access_token, refresh_token} = rawjson || {};
      if (access_token !== undefined && access_token !== null) {
        currentToken.value = access_token;
        currentRefreshToken.value = refresh_token;
        s3CurrentToken.value = access_token;
        toClipboard(access_token).finally(() => {
          ElMessage.success(`get access_token success: ${access_token}`);
        });
        Object.assign(requestInfo, request);
        Object.assign(responseInfo, response);
        Object.assign(rawJsonInfo, rawjson);
        Object.assign(exampleInfo, example);
        window.history.replaceState({}, document.title, window.location.pathname);
        updateReqAndRes();
      } else {
        ElMessage.error("Get token failed. Please refresh the user code!");
        resetData();
      }
    } else if (code === -1) {
      ElMessage.error('Authorization timeout. Please refresh the user code!')
      resetData();
    }
  }).catch((error) => {
    console.error('Error checking token:', error);
  })
}

function resetData() {
  verification_uri.value = "";
  user_code.value = "";
  device_code.value = "";
  expires_in.value = "";
  interval.value = "";
  currentToken.value = "";
  currentRefreshToken.value = "";
}

function handleRefreshToken() {
  if(props.cfgData.client_id.length === 0){
    ElMessage.error('client_id is empty, please click the config button on the right side, and check the configuration');
    return;
  }
  else if(props.cfgData.client_secret.length === 0){
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
        const {access_token, refresh_token} = rawjson || {};
        currentToken.value = access_token??"Uncertain";
        currentRefreshToken.value = refresh_token??"Uncertain";
        s3CurrentToken.value = access_token??"Uncertain";
        toClipboard(access_token).finally(() => {
          ElMessage.success(`refresh access_token success: ${access_token}`);
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
const contentType = ref("application/json");
const headerDialogVisible = ref(false);
const bodyDialogVisible = ref(false);
const newHeader = reactive({
  name: '',
  value: ''
});
const headerFormRef = ref<FormInstance>()
const headerRules = reactive({
  name: [{ required: true, message: 'Please enter header name', trigger: 'blur' }],
  value: [{ required: true, message: 'Please enter header value', trigger: 'blur' }]
});
const additionalHeaders = reactive([]);
const additionalBody = ref("")

const methodOptions = [
  {label: 'GET', value: 'GET', disabled: false},
  {label: 'POST', value: 'POST', disabled: false},
  {label: 'PUT', value: 'PUT', disabled: false},
  {label: 'DELETE', value: 'DELETE', disabled: false},
  {label: 'PATCH', value: 'PATCH', disabled: true},
  {label: 'HEAD', value: 'HEAD', disabled: true},
  {label: 'OPTIONS', value: 'OPTIONS', disabled: true}
];

const contentTypeOptions = [
  {label: 'application/json', value: 'application/json', disabled: false},
  {label: 'application/atom+xml', value: 'application/atom+xml', disabled: true},
  {label: 'text/plain', value: 'text/plain', disabled: true},
  {label: 'text/csv', value: 'text/csv', disabled: true},
  // {label: 'Custom', value: 'Custom', disabled: true}, // TODO：参考google playground
]

function handleMethodChange(value) {
  requestMethod.value = value;
  ElMessage.info(`当前请求方法: ${value}`);
}

function handleContentTypeChange(value) {
  contentType.value = value;
  ElMessage.info(`当前Content-Type: ${value}`);
}

async function handleRequestAPI() {
  if(requestUri.value.length === 0){
    ElMessage.error('api address cannot be empty');
    return;
  }else if(s3CurrentToken.value.length === 0){
    ElMessage.error('access_token is empty, please get the access_token firstly');
    return;
  }else{
    // 将 additionalHeaders 数组转换对象
    const additionalHeadersObject = additionalHeaders.reduce((map, header) => {
      if (header.name && header.value) {
        map[header.name] = header.value;
      }
      return map;
    }, {});
    const dataObject = {
      method: requestMethod.value,
      api_addr: requestUri.value,
      access_token: s3CurrentToken.value,
      access_token_type: s3TokenType.value,
      header: additionalHeadersObject,
      http_body: additionalBody.value
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

function isJsonResponse(header) {
  const contentType = header["Content-Type"];
  return contentType != null && contentType.includes("application/json");
}

function formatJson(jsonStr) {
  // 格式化 JSON 内容
  try {
    let resStr = JSON.stringify(JSON.parse(jsonStr), null, '  ');
    return resStr;
  } catch (error) {
    // 解析失败，返回原始内容
    console.error('格式化json失败');
    return jsonStr;
  }
}
async function addHeader(headerForm) {
  if (!headerForm)
    return;
  await headerForm.validate((valid, fields) => {
    if (valid) {
      additionalHeaders.push({name: newHeader.name, value: newHeader.value});
      // 添加之后清空
      newHeader.name = '';
      newHeader.value = '';
    } else {
      console.log('error add', fields)
    }
  })
}

function deleteRow(index) {
  additionalHeaders.splice(index, 1)
}

async function generateQRCode(url) {
  try {
    qrCodeSrc.value = await QRCode.toDataURL(url);
    console.log()
  } catch (err) {
    console.error('Failed to generate QR code: ', err)
  }
}

watch(props.cfgData, (newValue) => {
  s1Data.authorization_endpoint = newValue.authorization_endpoint;
  s1Data.scope = newValue.default_scope;
  initialAddress.value = newValue.authorization_endpoint.concat(
      "?response_type=device_code",
      newValue.default_scope?.length > 0 ? "&scope=".concat(newValue.default_scope) : "",
      newValue.client_id?.length > 0 ? "&client_id=".concat(newValue.client_id) : "",
      // "&redirect_uri=",
      // s1Data.redirect_uri,
      // s1Data.state?.length > 0 ? "&state=".concat(s1Data.state) : ""
  );
  requestUri.value = newValue.userinfo_endpoint;
  s3TokenType.value = newValue.access_token_type;
});

watch(verification_uri, (newValue) => {
  generateQRCode(newValue);
});

onMounted(() => {
  requestUri.value = props.cfgData.userinfo_endpoint;
  s3TokenType.value = props.cfgData.access_token_type;
  additionalHeaders.length = 0;
  additionalBody.value = '';

  if (verification_uri.value.length != 0) {
    generateQRCode(verification_uri.value);
  }
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
            <span class="stepTitle">Step 1: Request for Device Flow Authorization</span>
          </template>
          <el-scrollbar class="fitSide" ref="agS1ContainerRef">
            <h4 style="text-align: left;margin: 0">accessToken Endpoint</h4>
            <el-input v-model="s1Data.authorization_endpoint" disabled/>
<!--            <h4 style="text-align: left;margin: 0">Redirect Uri</h4>-->
<!--            <el-input v-model="s1Data.redirect_uri" disabled/>-->
            <h4 style="text-align: left;margin: 0">Scope</h4>
            <el-input v-model="s1Data.scope" placeholder="Scope" @input="handleS1Change" @blur="handleS1Change"/>
            <h4 style="text-align: left;margin: 0">Response Type</h4>
            <el-input v-model="s1Data.response_type" placeholder="Response Type" disabled/>
            <h4 style="text-align: left;margin: 0">Grant Url</h4>
            <el-input v-model="initialAddress" type="textarea" :autosize="{ minRows: 4, maxRows: 6 }"
                      placeholder="Request grant code address"/>
            <el-button ref="agS1Ref" type="primary" @click="handleDeviceFlow" class="n-button side-button">
              GO
            </el-button>
          </el-scrollbar>
        </el-collapse-item>
        <el-collapse-item class="el-collapse-item" name="2">
          <template #title>
            <span class="stepTitle">Step 2: Fetch token with Device Flow</span>
          </template>
          <el-scrollbar class="fitSide">
            <h4 style="text-align: left;margin: 0">User Code</h4>
            <span style="display: flex">The user code will expire in 5 minutes</span>
            <el-input v-model="user_code" placeholder="User Code">
              <template #append>
                <el-button type="primary" @click="handleDeviceFlow" class="t-button">
                  <span style="margin-left: -10px">Fresh User Code</span>
                </el-button>
              </template>
            </el-input>
            <h4 style="text-align: left;margin: 0; margin-top: 5px">Verification Uri</h4>
            <span style="display: flex">Go to the following link or scan the QR code,then enter the user code above</span>
            <el-input v-model="verification_uri" placeholder="Verification Uri" />
            <div>
              <img :src="qrCodeSrc" v-if="qrCodeSrc != ''" alt="QR Code" />
            </div>
            <h4 style="text-align: left;margin: 0">Refresh Token</h4>
            <el-input v-model="currentRefreshToken" placeholder="Refresh Token">
              <template #append>
                <el-button type="primary" @click="handleRefreshToken" class="t-button">
                  Refresh Token
                </el-button>
              </template>
            </el-input>
            <h4 style="text-align: left;margin: 0">Current Access Token</h4>
            <div v-if="currentToken" class="tokenArea">
              <strong><code>{{ currentToken }}</code></strong>
            </div>
            <div v-else class="tokenArea">
              <strong><code>Waiting for authorization...</code></strong>
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
            <h4 style="text-align: left;margin: 0">Method</h4>
            <el-select v-model="requestMethod" @change="handleMethodChange">
              <el-option v-for="item in methodOptions" :key="item.label" :label="item.label" :value="item.value"
                         :disabled="item.disabled"/>
            </el-select>
            <el-row>
              <el-col :span="6">
                <el-button
                    @click="headerDialogVisible = true"
                    class="user-button">
                  Add Headers
                </el-button>
              </el-col>
              <el-col :span="8">
                <el-button @click="bodyDialogVisible = true" class="user-button">
                  Enter Request Body
                </el-button>
              </el-col>
              <el-col :span="10" style="margin-top: 10px">
                <el-select v-model="contentType" @change="handleContentTypeChange">
                  <el-option v-for="item in contentTypeOptions" :key="item.label" :label="item.label" :value="item.value"
                             :disabled="item.disabled"/>
                </el-select>
              </el-col>
            </el-row>
            <!-- header Dialog-->
            <el-dialog v-model="headerDialogVisible" title="Headers">
              <div>
                <!--                <h4 style="display: flex">Add a header:</h4>-->
                <el-form ref="headerFormRef" :inline="true" :model="newHeader" :rules="headerRules" style="display: flex">
                  <el-form-item prop="name" style="width: 280px">
                    <el-input v-model="newHeader.name" placeholder="header name"></el-input>
                  </el-form-item>
                  <el-form-item prop="value" style="width: 280px">
                    <el-input v-model="newHeader.value" placeholder="header value"></el-input>
                  </el-form-item>
                  <el-form-item style="width: 260px">
                    <el-button @click="addHeader(headerFormRef)" style="margin-bottom: 10px; margin-left: -10px">Add</el-button>
                  </el-form-item>
                </el-form>
                <!-- 展示自定义的header-->
                <el-table :data="additionalHeaders" border style="width: 100%">
                  <el-table-column label="Header Name" width="260" > <!-- 一定要设置width，不然会报错-->
                    <template v-slot="{ row }">
                      {{ row.name }}
                    </template>
                  </el-table-column>
                  <el-table-column label="Header Value" width="260">
                    <template v-slot="{ row }">
                      {{ row.value }}
                    </template>
                  </el-table-column>
                  <el-table-column fixed="right" width="260">
                    <template #default="scope">
                      <el-button
                          type="primary"
                          size="small"
                          @click="deleteRow(scope.$index)"
                      >
                        Remove
                      </el-button>
                    </template>
                  </el-table-column>
                </el-table>
                <el-button @click="headerDialogVisible = false" style="display: flex">Close</el-button>
              </div>
            </el-dialog>
            <!-- header Dialog-->
            <el-dialog v-model="bodyDialogVisible" title="Request Body">
              <h4 style="display: flex">Enter the data that will be added to the body of the request:</h4>
              <el-input
                  v-model="additionalBody"
                  :rows="15"
                  type="textarea"
                  placeholder="Please input"
              />
              <el-button @click="bodyDialogVisible = false" style="display: flex">Close</el-button>
            </el-dialog>
            <h4 style="text-align: left;margin: 0">Access Token</h4>
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
          <!--   <el-checkbox style="position: absolute;bottom: 30px;left: 20px;" v-model="isWrapReq" label="Wrap Lines"
                          size="large"/> -->
        </div>
      </div>
      <div class="http-right">
        <el-divider content-position="left" direction="horizontal" class="http-info-type">
          <span>RESPONSE INFO</span>
        </el-divider>
        <div class="http-content" style="text-align: start; padding: 0em; position: relative; overflow: auto; max-height: 300px; width: 100%">
          <el-scrollbar class="http-render">
            <highlightjs autodetect :code="responseInfo.strHeader"/>
            <highlightjs v-if="isJsonResponse(responseInfo.header)" autodetect :code="formatJson(responseInfo.body)"/>
            <highlightjs v-else autodetect :code="responseInfo.body"></highlightjs>
          </el-scrollbar>
          <el-checkbox v-model="isWrapRes" label="Wrap Lines"
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

    .user-button {
      color: black;
      background-color: rgba(128, 128, 128, 0.2);
      width: 140px;
    }
    .user-button:hover {
      color: black;
      background-color: rgba(128, 128, 128, 0.4);
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
