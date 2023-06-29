const protocol = window.location.protocol;
const hostname = window.location.hostname;
const port = window.location.port;
const prefix = `${protocol}//${hostname}${port ? ':' + port : ''}`;

const evn = 0;
const devBase = "http://10.24.65.197";
const prodBase = prefix;

export function getApiBase() {
  let _l = "";
  switch(evn){
    // @ts-ignore
    case 0:
      _l = devBase;
      break;
    // @ts-ignore
    case 1:
      _l = prodBase;
      break;
    default:
      _l = devBase;
  }
  return _l;
}

export const getParam = (_k) => {
  const pv = new RegExp(`(^|&)${  _k  }=([^&]*)(&|$)`, "i");
  const r = window.location.search.slice(1).match(pv);
  if(r !== null){
    return decodeURIComponent(r[2]);
  }
  return '';
};

export const generateRandomString = (input: string) => {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  const timestamp = Date.now().toString();
  const randomLength = Math.floor(Math.random() * input.length) + 1;
  let result = '';
  for(let i = 0; i < randomLength; i++){
    result += input.charAt(Math.floor(Math.random() * input.length));
  }
  for(let i = 0; i < 14 - randomLength; i++){
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  result += timestamp.slice(-2);
  return result;
};
