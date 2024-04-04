<script setup>
import RegistrySearch from "./components/RegistrySearch.vue";
import SecurityMine from "./components/securityMine.vue";
import CveSearch from "./components/CveSearch.vue";
import Todo from "./components/todo.vue";
</script>

<template>
  <n-card title="Pikachu">
    <n-tabs type="bar" animated class="tabs">
      <n-tab-pane name="注册表扫描" tab="注册表扫描" display-directive="show:lazy" class="tabs">
        <RegistrySearch/>
      </n-tab-pane>
      <n-tab-pane name="安全自检" display-directive="show:lazy" tab="安全自检">
        <securityMine/>
      </n-tab-pane>
      <n-tab-pane name="CVE" display-directive="show:lazy" tab="CVE漏洞库">
        <CveSearch/>
      </n-tab-pane>
      <n-tab-pane name="todo" display-directive="show:lazy" tab="todo">
<!--        <todo/>-->123
      </n-tab-pane>
    </n-tabs>
  </n-card>
</template>

<style scoped>
.tabs{
  font-weight: bold;
  font-size: 20rpx;
}
</style>
