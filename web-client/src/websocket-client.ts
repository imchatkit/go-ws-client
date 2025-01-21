import { ConnectionState, ClientOptions, EventHandlers } from './types';

/**
 * WebSocket 客户端实现
 */
export class WebSocketClient {
    private ws: WebSocket | null = null;
    private state: ConnectionState = ConnectionState.Disconnected;
    private options: ClientOptions;
    private reconnectTimer: number | null = null;
    private heartbeatTimer: number | null = null;
    private handlers: EventHandlers = {};

    /**
     * 构造函数
     * @param options 客户端配置选项
     */
    constructor(options: Partial<ClientOptions> = {}) {
        const defaultOptions: ClientOptions = {
            autoReconnect: false,
            reconnectInterval: 5000,
            heartbeatInterval: 30000,
            debug: false
        };
        this.options = { ...defaultOptions, ...options };
    }

    /**
     * 建立 WebSocket 连接
     * @param url WebSocket 服务器地址
     * @param token 认证令牌
     * @param deviceType 设备类型标识
     * @returns 是否成功发起连接
     */
    connect(url: string, token: string, deviceType: string = 'im_app_android'): boolean {
        if (this.state === ConnectionState.Connected) {
            this.debug('Already connected');
            return true;
        }

        try {
            this.setState(ConnectionState.Connecting);

            // 添加认证参数到 URL
            const wsUrl = new URL(url);
            wsUrl.searchParams.set('token', token);
            wsUrl.searchParams.set('deviceType', deviceType);

            // 创建 WebSocket 连接
            this.ws = new WebSocket(wsUrl.toString());
            
            // 设置事件处理器
            this.setupEventHandlers();
            
            return true;
        } catch (error) {
            this.debug('Connection failed:', error);
            this.setState(ConnectionState.Disconnected);
            return false;
        }
    }

    /**
     * 设置事件处理器
     */
    private setupEventHandlers(): void {
        if (!this.ws) return;

        this.ws.onopen = () => {
            this.setState(ConnectionState.Connected);
            this.debug('WebSocket connected');
            this.startHeartbeat();
        };

        this.ws.onmessage = (event: MessageEvent) => {
            this.debug('Received message:', event.data);
            if (this.handlers.onMessage) {
                this.handlers.onMessage(event.data);
            }
        };

        this.ws.onerror = (error: Event) => {
            this.debug('WebSocket error:', error);
            if (this.handlers.onError) {
                this.handlers.onError(error);
            }
        };

        this.ws.onclose = () => {
            this.setState(ConnectionState.Disconnected);
            this.debug('WebSocket closed');
            this.stopHeartbeat();
            
            if (this.handlers.onClose) {
                this.handlers.onClose();
            }

            if (this.options.autoReconnect) {
                this.scheduleReconnect();
            }
        };
    }

    /**
     * 发送消息到服务器
     * @param message 要发送的消息
     * @returns 是否成功发送
     */
    send(message: string): boolean {
        if (!this.isConnected()) {
            this.debug('Not connected');
            return false;
        }

        try {
            this.ws?.send(message);
            return true;
        } catch (error) {
            this.debug('Send failed:', error);
            return false;
        }
    }

    /**
     * 关闭 WebSocket 连接
     * @returns 是否成功关闭
     */
    close(): boolean {
        if (!this.isConnected()) {
            return false;
        }

        try {
            this.options.autoReconnect = false; // 禁用自动重连
            this.ws?.close();
            this.setState(ConnectionState.Disconnected);
            this.stopHeartbeat();
            this.cancelReconnect();
            return true;
        } catch (error) {
            this.debug('Close failed:', error);
            return false;
        }
    }

    /**
     * 设置事件处理器
     * @param handlers 事件处理器对象
     */
    setEventHandlers(handlers: EventHandlers): void {
        this.handlers = { ...this.handlers, ...handlers };
    }

    /**
     * 检查是否已连接
     * @returns 是否处于已连接状态
     */
    isConnected(): boolean {
        return this.state === ConnectionState.Connected;
    }

    /**
     * 获取当前连接状态
     * @returns 当前状态
     */
    getState(): ConnectionState {
        return this.state;
    }

    /**
     * 更新连接状态
     * @param state 新状态
     */
    private setState(state: ConnectionState): void {
        this.state = state;
        if (this.handlers.onStateChange) {
            this.handlers.onStateChange(state);
        }
    }

    /**
     * 开始心跳
     */
    private startHeartbeat(): void {
        this.stopHeartbeat();
        this.heartbeatTimer = window.setInterval(() => {
            this.send('ping');
        }, this.options.heartbeatInterval);
    }

    /**
     * 停止心跳
     */
    private stopHeartbeat(): void {
        if (this.heartbeatTimer) {
            clearInterval(this.heartbeatTimer);
            this.heartbeatTimer = null;
        }
    }

    /**
     * 安排重连
     */
    private scheduleReconnect(): void {
        this.cancelReconnect();
        this.setState(ConnectionState.Reconnecting);
        this.reconnectTimer = window.setTimeout(() => {
            this.debug('Attempting to reconnect...');
            this.connect(this.ws?.url || '', '', '');
        }, this.options.reconnectInterval);
    }

    /**
     * 取消重连
     */
    private cancelReconnect(): void {
        if (this.reconnectTimer) {
            clearTimeout(this.reconnectTimer);
            this.reconnectTimer = null;
        }
    }

    /**
     * 输出调试信息
     * @param message 调试消息
     * @param args 其他参数
     */
    private debug(message: string, ...args: any[]): void {
        if (this.options.debug) {
            console.log(`[WebSocketClient] ${message}`, ...args);
        }
    }
} 