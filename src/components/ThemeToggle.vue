<template>
  <n-tooltip trigger="hover">
    <template #trigger>
      <n-button
        :text="!showLabel"
        :aria-label="tooltipText"
        @click="cycleMode"
      >
        <template #icon>
          <n-icon :size="showLabel ? 18 : 20">
            <component :is="modeIcon" />
          </n-icon>
        </template>
        <template v-if="showLabel">{{ MODE_LABELS[mode] }}</template>
      </n-button>
    </template>
    {{ tooltipText }}
  </n-tooltip>
</template>

<script setup>
import { computed } from "vue";
import { SunnyOutline, MoonOutline, DesktopOutline } from "@vicons/ionicons5";
import { useTheme, MODE_LABELS } from "../composables/useTheme";

defineProps({
  showLabel: Boolean,
});

const { mode, nextMode, cycleMode } = useTheme();

const MODE_ICONS = {
  system: DesktopOutline,
  light: SunnyOutline,
  dark: MoonOutline,
};

const modeIcon = computed(() => MODE_ICONS[mode.value]);

const tooltipText = computed(
  () =>
    `Theme: ${MODE_LABELS[mode.value]} (click for ${MODE_LABELS[nextMode.value]})`,
);
</script>
