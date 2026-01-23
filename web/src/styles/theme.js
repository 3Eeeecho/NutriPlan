/**
 * Naive UI 主题配置
 * 基于 NutriPlan 设计系统定制
 */

import { colors } from './tokens'

export const themeOverrides = {
  common: {
    // 主色调
    primaryColor: colors.primary.default,
    primaryColorHover: colors.primary.hover,
    primaryColorPressed: colors.primary.pressed,
    primaryColorSuppl: colors.primary.light,

    // 信息色
    infoColor: colors.info.default,
    infoColorHover: '#2563eb',
    infoColorPressed: '#1d4ed8',
    infoColorSuppl: colors.info.light,

    // 成功色
    successColor: colors.success.default,
    successColorHover: '#059669',
    successColorPressed: '#047857',
    successColorSuppl: colors.success.light,

    // 警告色
    warningColor: colors.warning.default,
    warningColorHover: '#d97706',
    warningColorPressed: '#b45309',
    warningColorSuppl: colors.warning.light,

    // 错误色
    errorColor: colors.error.default,
    errorColorHover: '#dc2626',
    errorColorPressed: '#b91c1c',
    errorColorSuppl: colors.error.light,

    // 文字颜色
    textColorBase: colors.text.primary,
    textColor1: colors.text.primary,
    textColor2: colors.text.secondary,
    textColor3: colors.text.tertiary,
    textColorDisabled: colors.text.disabled,

    // 边框颜色
    borderColor: colors.border.primary,
    dividerColor: colors.border.secondary,

    // 圆角
    borderRadius: '8px',
    borderRadiusSmall: '6px',

    // 字体
    fontFamily:
      '-apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, "Helvetica Neue", Arial, "Noto Sans", sans-serif',
    fontSize: '14px',
    fontSizeMini: '12px',
    fontSizeTiny: '12px',
    fontSizeSmall: '13px',
    fontSizeMedium: '14px',
    fontSizeLarge: '15px',
    fontSizeHuge: '16px',

    // 高度
    heightTiny: '28px',
    heightSmall: '32px',
    heightMedium: '36px',
    heightLarge: '40px',
    heightHuge: '44px',
  },

  // 按钮
  Button: {
    borderRadiusMedium: '8px',
    borderRadiusLarge: '10px',
    fontSizeMedium: '14px',
    fontSizeLarge: '15px',
    fontWeightStrong: '500',
    paddingMedium: '0 18px',
    paddingLarge: '0 22px',
    textColorPrimary: '#ffffff',
  },

  // 卡片
  Card: {
    borderRadius: '12px',
    paddingMedium: '20px',
    paddingLarge: '24px',
    borderColor: colors.border.primary,
    color: '#ffffff',
  },

  // 输入框
  Input: {
    borderRadius: '8px',
    borderHover: colors.primary.light,
    borderFocus: colors.primary.default,
    heightMedium: '36px',
    heightLarge: '40px',
  },

  // 选择器
  Select: {
    peers: {
      InternalSelection: {
        borderRadius: '8px',
        heightMedium: '36px',
        heightLarge: '40px',
      },
    },
  },

  // 对话框
  Dialog: {
    borderRadius: '16px',
  },

  // 抽屉
  Drawer: {
    borderRadius: '16px 0 0 16px',
  },

  // 消息提示
  Message: {
    borderRadius: '8px',
  },

  // 标签
  Tag: {
    borderRadius: '6px',
  },

  // 进度条
  Progress: {
    fillColor: colors.primary.default,
    railColor: colors.bg.tertiary,
  },

  // 数据表格
  DataTable: {
    borderRadius: '8px',
    borderColor: colors.border.primary,
  },

  // 分页
  Pagination: {
    itemBorderRadius: '6px',
  },
}

// 深色主题覆盖（预留）
export const darkThemeOverrides = {
  common: {
    ...themeOverrides.common,
    bodyColor: '#0f172a',
    cardColor: '#1e293b',
    modalColor: '#1e293b',
    popoverColor: '#1e293b',
    textColorBase: '#f8fafc',
    textColor1: '#f8fafc',
    textColor2: '#cbd5e1',
    textColor3: '#94a3b8',
    borderColor: '#334155',
    dividerColor: '#475569',
  },
}
