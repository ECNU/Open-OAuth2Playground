<script setup lang="ts">
import { onMounted, reactive, ref, watch } from 'vue'
import {ElMessage, FormInstance} from 'element-plus'
import useClipboard from 'vue-clipboard3';
import { fetchACTokenByClient, fetchApiData } from "/@/api/playground";

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
const currentToken = ref("");

async function handleGetTokenByClient() {
  if(props.cfgData.client_id.length === 0){
    ElMessage.error('client_id cannot be empty');
    return;
  }else if(props.cfgData.client_secret.length === 0){
    ElMessage.error('client_secret cannot be empty');
    return;
  }else{
    const dataObject = {
      client_id: props.cfgData.client_id,
      client_secret: props.cfgData.client_secret
    };
    fetchACTokenByClient(dataObject).then(({code, msg, data}) => {
      if(code === 0){
        const {request, response, rawjson, example} = data;
        const {access_token} = rawjson || {};
        currentToken.value = access_token??"Uncertain";
        s3CurrentToken.value = access_token??"Uncertain";
        toClipboard(access_token).finally(() => {
          ElMessage.success(`get access_token success: ${access_token}`);
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
  // {label: 'Custom', value: 'Custom', disabled: true}, // TODO：参考google，跳出添加header的dialog，header name = “Content-Type”,value由用户填
]

function handleMethodChange(value) {
  requestMethod.value = value;
  ElMessage.info(`当前请求方法: ${requestMethod.value}`);
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
    // console.log(resStr);
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

watch(props.cfgData, (newValue) => {
  requestUri.value = newValue.userinfo_endpoint;
  s3TokenType.value = newValue.access_token_type;
});

onMounted(() => {
  requestUri.value = props.cfgData.userinfo_endpoint;
  s3TokenType.value = props.cfgData.access_token_type;
  additionalHeaders.length = 0;
  additionalBody.value = '';
});


</script>
<template>
  <el-col :xs="24" :sm="24" :md="8" :lg="8" :xl="8">
    <div class="demo-collapse" id="main">
      <el-collapse v-model="activeName" accordion>
        <el-collapse-item class="el-collapse-item" name="1">
          <template #title>
            <span class="stepTitle">Step 1: Fetch token with Client Credentials</span>
          </template>
          <el-scrollbar class="fitSide">
            <el-button type="primary" @click="handleGetTokenByClient" class="t-button">
              Get Token
            </el-button>

            <h4 style="text-align: left;margin: 0">Current Access Token</h4>
            <div class="tokenArea">
              <strong><code>{{ currentToken }}</code></strong>
            </div>
          </el-scrollbar>
        </el-collapse-item>
        <el-collapse-item class="el-collapse-item" name="3">
          <template #title>
            <span class="stepTitle">Step 2: Request to API with the <code
                style="color:#cd3221">access_token</code></span>
          </template>
          <el-scrollbar class="fitSide">
            <h4 style="text-align: left;margin: 0">Request URI</h4>
            <el-input v-model="requestUri" placeholder="Authorization Code"/>
            <h4 style="text-align: left;margin: 0">Method</h4>
            <el-select v-model="requestMethod" @change="handleMethodChange" >
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
        <div class="http-content" style="text-align: start; padding: 0em; position: relative; overflow: auto; max-height: 350px; width: 100%">
          <el-scrollbar class="http-render">
            <highlightjs autodetect :code="responseInfo.strHeader"/>
            <!-- TODO: json区增加复制-->
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
