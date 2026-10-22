import { ref, computed, watch } from "vue";
import { darkTheme } from "naive-ui";

// Kept outside the "oneoff_cache_" prefix so cache clears don't reset it.
// index.html reads the same key before mount to avoid a flash.
const STORAGE_KEY = "oneoff_theme";
const MODES = ["system", "light", "dark"];
const THEME_COLORS = { dark: "#18181c", light: "#ffffff" };

export const MODE_LABELS = {
  system: "System",
  light: "Light",
  dark: "Dark",
};

const readMode = () => {
  try {
    const stored = localStorage.getItem(STORAGE_KEY);
    return MODES.includes(stored) ? stored : "system";
  } catch {
    return "system";
  }
};

const writeMode = (value) => {
  try {
    localStorage.setItem(STORAGE_KEY, value);
  } catch (error) {
    console.warn("Failed to save theme preference:", error);
  }
};

// Module-level singleton shared by every useTheme() caller
const mode = ref(readMode());
const systemPrefersDark = ref(false);

if (typeof window !== "undefined" && window.matchMedia) {
  const query = window.matchMedia("(prefers-color-scheme: dark)");
  systemPrefersDark.value = query.matches;
  query.addEventListener("change", (event) => {
    systemPrefersDark.value = event.matches;
  });
}

const resolved = computed(() => {
  if (mode.value === "system") {
    return systemPrefersDark.value ? "dark" : "light";
  }
  return mode.value;
});

// Naive UI treats a null theme as light
const naiveTheme = computed(() =>
  resolved.value === "dark" ? darkTheme : null,
);

const nextMode = computed(
  () => MODES[(MODES.indexOf(mode.value) + 1) % MODES.length],
);

const applyToDocument = (value) => {
  if (typeof document === "undefined") return;
  const root = document.documentElement;
  root.dataset.theme = value;
  root.style.colorScheme = value;
  const meta = document.querySelector('meta[name="theme-color"]');
  if (meta) meta.setAttribute("content", THEME_COLORS[value]);
};

watch(resolved, applyToDocument, { immediate: true });

const setMode = (value) => {
  if (!MODES.includes(value)) return;
  mode.value = value;
  writeMode(value);
};

const cycleMode = () => setMode(nextMode.value);

export function useTheme() {
  return {
    mode,
    nextMode,
    resolved,
    naiveTheme,
    setMode,
    cycleMode,
  };
}
