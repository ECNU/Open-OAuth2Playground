<script setup lang="ts">
import { ref } from "vue";
import Playground from "./playground";

defineOptions({
name: "Layout"
});

const activePage = ref("1");

function handlePageChange(index: string) {
  activePage.value = index;
}

function isPageActive(index: string) {
  return activePage.value === index;
}

</script>

<template>
  <el-container>
    <el-header>
      <el-menu default-active="1" class="el-menu-demo" mode="horizontal" @select="handlePageChange">
        <div class="logo">LOGO</div>
        <el-menu-item index="1">PLAYGROUND</el-menu-item>
        <el-menu-item index="2">MENU1</el-menu-item>
        <el-menu-item index="3">MENU2</el-menu-item>
      </el-menu>
      <Playground v-if="isPageActive('1')"/>
    </el-header>
  </el-container>
</template>
<style lang="less" scoped>
.el-header {
  padding: 0;

  .el-menu {
    background: linear-gradient(#fafafa 0%, #f5f5f5 100%);
    --el-menu-active-color: #666;
    --el-menu-text-color: var(--el-text-color-primary);
    --el-menu-hover-text-color: #333;
    --el-menu-bg-color: #eee;
    --el-menu-hover-bg-color: #eee;
    --el-menu-item-height: 56px;
    --el-menu-item-font-size: var(--el-font-size-base);
    --el-menu-item-hover-fill: var(--el-color-primary-light-9);
    --el-menu-border-color: var(--el-border-color);

    .logo {
      color: #fff;
      font-size: 20px;
      font-weight: bold;
      line-height: 56px;
      padding: 0 20px;
      background: #b70031;
    }
  }
}

</style>

