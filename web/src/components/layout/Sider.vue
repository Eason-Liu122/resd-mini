<template>
  <div class="flex pb-2 flex-col h-full min-w-[80px] border-r border-slate-100 dark:border-slate-800">
    <div class="w-full flex flex-row items-center justify-center pt-5 ml-[-5px]">
      <img class="w-12 h-12 cursor-pointer" src="@/assets/image/logo.png"  alt="resd-mini logo"/>
    </div>
    <main class="flex-1 flex-grow-1 mb-5 overflow-auto flex flex-col pt-1 items-center h-full">
      <NScrollbar :size="1">
        <NLayout has-sider>
          <NLayoutSider
              :bordered="false"
              show-trigger
              collapse-mode="width"
              :on-after-enter="() => { showAppName = true }"
              :on-after-leave="() => { showAppName = false }"
              :collapsed-width="70"
              :default-collapsed="true"
              :width="120"
              :native-scrollbar="false"
              :inverted="inverted"
              class="bg-inherit"
          >
            <NMenu
                :inverted="inverted"
                :collapsed-width="70"
                :collapsed-icon-size="22"
                :options="menuOptions"
                :value="menuValue"
                @update:value="handleUpdateValue"
            />
          </NLayoutSider>
        </NLayout>
        <NLayoutFooter position="absolute" :inverted="inverted" class="bg-inherit">
          <NMenu
              :inverted="inverted"
              :collapsed-width="70"
              :collapsed-icon-size="22"
              :options="footerOptions"
              :value="menuValue"
              @update:value="handleFooterUpdate"
          />
        </NLayoutFooter>
      </NScrollbar>
    </main>
  </div>
  <Footer v-model:showModal="showAppInfo" />
</template>

<script lang="ts" setup>
import type {MenuOption} from "naive-ui"
import {NIcon} from "naive-ui"
import {computed, h, ref, watch} from "vue"
import {useRoute, useRouter} from "vue-router"
import {
  CloudOutline,
  SettingsOutline,
  HelpCircleOutline,
  MoonOutline, SunnyOutline, LogoGithub
} from "@vicons/ionicons5"
import {useIndexStore} from "@/stores"
import Footer from "@/components/Footer.vue"

const route = useRoute()
const router = useRouter()
const inverted = ref(false)
const showAppName = ref(false)
const showAppInfo = ref(false)
const menuValue = ref(route.fullPath.substring(1))
const store = useIndexStore()

const globalConfig = computed(()=>{
  return store.globalConfig
})

const theme = computed(() => {
  return store.globalConfig.Theme === "darkTheme" ? renderIcon(SunnyOutline) : renderIcon(MoonOutline)
})

watch(() => route.path, (newPath, oldPath) => {
  menuValue.value = route.fullPath.substring(1)
});

const renderIcon = (icon: any) => {
  return () => h(NIcon, null, {default: () => h(icon)})
}

const menuOptions = ref([
  {
    label: "拦截",
    key: 'index',
    icon: renderIcon(CloudOutline),
  },
  {
    label: "设置",
    key: 'setting',
    icon: renderIcon(SettingsOutline),
  },
])

const footerOptions = ref([
  {
    label: "主题",
    key: 'theme',
    icon: theme,
  },
  {
    label: "github",
    key: 'github',
    icon: renderIcon(LogoGithub),
  },
  {
    label: "关于",
    key: 'about',
    icon: renderIcon(HelpCircleOutline),
  },
])

const handleUpdateValue = (key: string, item: MenuOption) => {
  menuValue.value = key
  return router.push({path: "/" + key})
}
const handleFooterUpdate = (key: string, item: MenuOption) => {
  if (key === "about") {
    showAppInfo.value = true
    return
  }

  if (key === "github") {
    window.open("https://github.com/putyy/resd-mini", "_blank");
    return
  }
  if (key === "theme") {
    if (globalConfig.value.Theme === "darkTheme") {
      store.setConfig(Object.assign({}, globalConfig.value, {Theme: "lightTheme"}))
      return
    }
    store.setConfig(Object.assign({}, globalConfig.value, {Theme: "darkTheme"}))

    return
  }
  menuValue.value = key
  return router.push({path: "/" + key})
}

</script>