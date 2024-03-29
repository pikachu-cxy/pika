<script setup>
import {SearchRegistry} from "../../wailsjs/go/main/App";
import {EventsOn} from "../../wailsjs/runtime";
import {ref} from 'vue';
const data = ref([]);

const input = ref('');

const columns = ref([
  {
    title: '软件名称',
    key: 'key'
  },
  {
    title: '注册表路径',
    key: 'path'
  }
])

async function initRegistryMap(){
  console.log(input.value)
  if(input.value.length === 0){
    //input.value = ""
    return
  }
  await SearchRegistry(input.value)
  console.log(data.value)

}

EventsOn("SearchRegistry", e => {
  data.value.push({
    Key: e.key,
    Path: e.Path
  });
});

</script>

<template>
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
  <n-data-table
      size="small"
      :columns="columns"
      :data="data"
      :max-height="300"
      style="margin-top: 10px"
      striped
  >
  </n-data-table>
</template>

<style scoped>

</style>