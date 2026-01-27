/**
 * NutriPlan 设计系统 - JavaScript Token
 * 用于 Naive UI 主题配置和 JS 中使用
 */

export const colors = {
  // 品牌色
  primary: {
    default: "#059669", // Emerald 600
    hover: "#047857", // Emerald 700
    pressed: "#065f46", // Emerald 800
    light: "#34d399", // Emerald 400
    lighter: "#a7f3d0", // Emerald 200
    supperlight: "#ecfdf5", // Emerald 50
  },

  // 辅助色
  accent: {
    warm: "#f59e0b", // Amber 500
    warmLight: "#fef3c7", // Amber 100
    cool: "#3b82f6", // Blue 500
    coolLight: "#dbeafe", // Blue 100
    pink: "#ec4899", // Pink 500
    pinkLight: "#fce7f3", // Pink 100
    purple: "#8b5cf6", // Violet 500
    purpleLight: "#ede9fe", // Violet 100
  },

  // 功能色
  success: {
    default: "#10b981",
    light: "#d1fae5",
  },
  warning: {
    default: "#f59e0b",
    light: "#fef3c7",
  },
  error: {
    default: "#ef4444",
    light: "#fee2e2",
  },
  info: {
    default: "#3b82f6",
    light: "#dbeafe",
  },

  // 中性色
  text: {
    primary: "#111827",
    secondary: "#4b5563",
    tertiary: "#9ca3af",
    disabled: "#d1d5db",
    inverse: "#ffffff",
  },
  bg: {
    primary: "#ffffff",
    glass: "rgba(255, 255, 255, 0.8)",
    secondary: "#f8fafc",
    tertiary: "#f1f5f9",
    hover: "#f1f5f9",
  },
  border: {
    primary: "#e2e8f0",
    secondary: "#f1f5f9",
  },

  // 营养素色彩
  nutrition: {
    protein: "#ef4444",
    carb: "#f59e0b",
    fat: "#8b5cf6",
    fiber: "#10b981",
  },
};

export const spacing = {
  xs: "4px",
  sm: "8px",
  md: "12px",
  lg: "16px",
  xl: "24px",
  "2xl": "32px",
  "3xl": "48px",
  "4xl": "64px",
};

export const radius = {
  xs: "4px",
  sm: "6px",
  md: "8px",
  lg: "12px",
  xl: "16px",
  "2xl": "20px",
  full: "9999px",
};

export const fontSize = {
  xs: "12px",
  sm: "13px",
  base: "14px",
  md: "15px",
  lg: "16px",
  xl: "18px",
  "2xl": "20px",
  "3xl": "24px",
  "4xl": "32px",
  "5xl": "40px",
};

export const fontWeight = {
  normal: 400,
  medium: 500,
  semibold: 600,
  bold: 700,
};

export const breakpoints = {
  xs: 480,
  sm: 640,
  md: 768,
  lg: 1024,
  xl: 1280,
  "2xl": 1536,
};

export const shadows = {
  xs: "0 1px 2px rgba(0, 0, 0, 0.05)",
  sm: "0 1px 3px rgba(0, 0, 0, 0.1), 0 1px 2px rgba(0, 0, 0, 0.06)",
  md: "0 4px 6px rgba(0, 0, 0, 0.07), 0 2px 4px rgba(0, 0, 0, 0.05)",
  lg: "0 10px 15px rgba(0, 0, 0, 0.1), 0 4px 6px rgba(0, 0, 0, 0.05)",
  xl: "0 20px 25px rgba(0, 0, 0, 0.1), 0 10px 10px rgba(0, 0, 0, 0.04)",
  "2xl": "0 25px 50px rgba(0, 0, 0, 0.15)",
};

export const zIndex = {
  dropdown: 1000,
  sticky: 1020,
  fixed: 1030,
  modalBackdrop: 1040,
  modal: 1050,
  popover: 1060,
  tooltip: 1070,
};
