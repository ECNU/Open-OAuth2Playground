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
const isHorizontal = ref(true);
const toggleLayout = () => {
  isHorizontal.value = !isHorizontal.value;
};

// Http Info
const requestInfo = reactive({});
const responseInfo = reactive({});
const oauthInfo = reactive({});
const exampleInfo = reactive({});

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
    fetchACToken(dataObject).then(({code,msg,data}) => {
      if(code===0){
        const {request, response, oauth2, example} = data;
        const {access_token, refresh_token} = oauth2;
        currentToken.value = access_token;
        currentRefreshToken.value = refresh_token;
        s3CurrentToken.value = access_token;
        toClipboard(access_token).finally(() => {
          ElMessage.success(`已复制access_token: ${access_token}`);
        });
        Object.assign(requestInfo, request);
        Object.assign(responseInfo, response);
        Object.assign(oauthInfo, oauth2);
        Object.assign(exampleInfo, example);

        window.history.replaceState({}, document.title, window.location.pathname);
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

    fetchRefreshToken(dataObject).then(({code,msg,data}) => {
      if(code === 0){
        const {request, response, oauth2, example} = data;
        const {access_token, refresh_token} = oauth2;
        currentToken.value = access_token;
        currentRefreshToken.value = refresh_token;
        s3CurrentToken.value = access_token;
        toClipboard(access_token).finally(() => {
          ElMessage.success(`已复制新access_token: ${access_token}`);
        });
        Object.assign(requestInfo, request);
        Object.assign(responseInfo, response);
        Object.assign(oauthInfo, oauth2);
        Object.assign(exampleInfo, example);
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
    fetchApiData(dataObject).then(({code,msg,data}) => {
      if(code === 0){
        const {request, response, oauth2, example} = data;
        Object.assign(requestInfo, request);
        Object.assign(responseInfo, response);
        Object.assign(oauthInfo, oauth2);
        Object.assign(exampleInfo, example);
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
  // todo: 浮动球拖动效果考虑的因素比较多，bug点也比较多，后期有时间看看能否完善
  // todo: 1. 需要同时考虑手机和pc页面不同的场景元素移动计算有所区别
  // todo: 2. 拖动过程中，浮动球的位置需要跟随鼠标移动，但是浮动球的位置需要相对在容器内，就需要做边缘检测
  // todo: 3. 拖动的整个周期需要阻止单击事件流
  // todo: 4. 拖动刚开始需要检测用户是否是真的想要拖动，可以检测拖动的时间或距离来综合判断
  // todo: 5. pc端还要考虑浏览器的的选择事件
  // todo: 6. 边缘吸附的时候需要考虑浮动球的位置，是吸附在容器的边缘还是容器内的边缘

  let isDragging = false;
  let mouseOffset = {x: 0, y: 0};
  let touchOffset = {x: 0, y: 0};

  floatButton.addEventListener('mousedown', function (event) {

  });

  floatButton.addEventListener('touchstart', function (event) {

  });

  document.addEventListener('mousemove', function (event) {
    if(isDragging){

    }else{
    }
  });

  document.addEventListener('touchmove', function (event) {
    if(isDragging){

    }else{
    }
  });

  document.addEventListener('mouseup', function () {
    isDragging = false;
  });

  document.addEventListener('touchend', function () {
    isDragging = false;
  });
};

</script>
<template>
  <div style="text-align: center;width:100%">TODO</div>
</template>

<style scoped lang="less">
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
}
</style>
