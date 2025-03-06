<template>
  <NConfigProvider class="h-full" :theme="theme" :locale="zhCN">
    <NaiveProvider>
      <RouterView />
    </NaiveProvider>
    <NGlobalStyle />
    <NModalProvider />
  </NConfigProvider>
</template>

<script setup lang="ts">
import NaiveProvider from '@/components/NaiveProvider.vue'
import {darkTheme, lightTheme, zhCN} from 'naive-ui'
import {useIndexStore} from "@/stores"
import {useWsStore} from "@/stores/ws"
import {computed, onMounted, watch} from "vue"
import type {wsType} from "@/types/ws"

const store = useIndexStore()
const wsStore = useWsStore()

const theme = computed(() => {
  if (store.globalConfig.Theme === "darkTheme") {
    document.documentElement.classList.add('dark');
    return darkTheme
  }
  document.documentElement.classList.remove('dark');
  return lightTheme
})

onMounted(async () => {
  await store.init()
  wsStore.setPort(store.globalConfig.Port)
  wsStore.websocketInit()
  wsStore.bindMessageHandle({
    type: "message",
    event: (res: wsType.Message)=>{
      switch (res?.code) {
        case 0:
          window?.$message?.error(res.message)
          break
        case 1:
          window?.$message?.success(res.message)
          break
      }
    }
  })
  wsStore.bindMessageHandle({
    type: "clipboard",
    event: (res: wsType.Clipboard)=>{
      navigator.clipboard.writeText(res.content);
    }
  })

  wsStore.bindMessageHandle({
    type: "updateProxyStatus",
    event: (res: any)=>{
      store.updateProxyStatus(res)
    }
  })
})
</script>

<style scoped>
</style>
