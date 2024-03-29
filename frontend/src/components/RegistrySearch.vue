<script setup>
import {ref} from "vue"
import {Greet2} from '../../wailsjs/go/main/App'

const input = ref('');
const value = ref('');


async function initRegistryMap(){
  if(input.value.length === 0){
    value.value = "请输入软件名称！"
    //input.value = ""
    return
  }
  value.value = "程序正在扫描中，请稍等片刻~"
  value.value = (await Greet2(input.value)).join(" ")
}

const handleChange = (event) => {
  //console.log(event)
  const files = event.target.files;
  if (files.length > 0) {
    console.log(event)
    const file = files[0];
    readFileContent(file);
  }
};
const readFileContent = async (file) => {
  try {
    const content = await readFile(file);
    input.value = content;
  } catch (error) {
    console.error("Error reading the file:", error);
  }
};

const readFile = (file) => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      resolve(reader.result);
    };
    reader.onerror = () => {
      reject(reader.error);
    };
    reader.readAsText(file);
  });
};

</script>

<template>
  <div class="box">
    <n-input v-model:value="input" type="textarea"
             placeholder="输入关键字"
             :autosize="{
                      minRows: 5,
                      maxRows: 20
                    }"
             class="input-l"/>
    <n-button type="primary" class="button" @click="initRegistryMap">
      点击搜索
    </n-button>
    <div class="custom-file-upload">
      <label for="file-upload" class="custom-file-upload-label">
        <span style="font-family: inherit;">选择文件</span>
      </label>
      <input type="file" id="file-upload" @change="handleChange">
    </div>

    <n-input
        v-model:value="value"
        type="textarea"
        placeholder="结果显示"
        :autosize="{
        minRows: 10,
        maxRows: 20
      }"
        class="input"
    />
  </div>
</template>

<style scoped>
.box {
  width: 90vw;

  .input-l {
    float: left;
    text-align: left;
    width: 50%;
  }

  .button {
    margin-left: 10px;
    float: left;
  }

  .input {
    font-weight: bold;
    margin-top: 20px;
    float: left;
    text-align: left;
    width: 100%;
  }

  .custom-file-upload {
    padding-left: 20px;
    float: left;
    left: 20px;
    display: inline-block;
  }
}

.custom-file-upload-label {
  border-radius: 5px;
  display: inline-block;
  background-color: #18a058; /* 半透明背景颜色 */
  color: #fff;
  padding: 5px 10px;
  border: 1px solid #ccc; /* 灰色边框 */
  width: 60px;
  height: 25px;
  text-align: center;
  cursor: pointer;
}

/* 隐藏原生文件输入框 */
#file-upload {
  display: none;
}

.result {
  height: 20px;
  line-height: 20px;
  margin: 1.5rem auto;
}

.input-box {
  text-align: center;
}

</style>
