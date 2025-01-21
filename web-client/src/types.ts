/**
 * WebSocket 连接状态枚举
 */
export enum ConnectionState {
    Disconnected = 0,  // 连接已断开
    Connecting = 1,    // 正在连接中
    Connected = 2,     // 已连接
    Reconnecting = 3   // 正在重新连接
}

/**
 * WebSocket 客户端配置选项
 */
export interface ClientOptions {
    autoReconnect: boolean;      // 是否自动重连
    reconnectInterval: number;   // 重连间隔时间（毫秒）
    heartbeatInterval: number;   // 心跳间隔时间（毫秒）
    debug: boolean;             // 是否启用调试模式
}

/**
 * 事件处理器类型定义
 */
export interface EventHandlers {
    onMessage?: (message: string) => void;           // 消息接收回调
    onStateChange?: (state: ConnectionState) => void; // 状态变更回调
    onError?: (error: Event) => void;                // 错误处理回调
    onClose?: () => void;                            // 连接关闭回调
} 