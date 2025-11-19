import {Activity, LayoutGrid, List, LogIn, Server, Sparkles} from 'lucide-react';
import GithubSvg from "../assets/github.svg";

interface PublicHeaderProps {
    title: string;
    lastUpdated?: string;
    viewMode?: 'grid' | 'list';
    onViewModeChange?: (mode: 'grid' | 'list') => void;
    showViewToggle?: boolean;
}

const PublicHeader = ({
                          title,
                          lastUpdated,
                          viewMode,
                          onViewModeChange,
                          showViewToggle = false
                      }: PublicHeaderProps) => {
    return (
        <header className="sticky top-0 z-50 border-b border-slate-200/80 bg-white/80 backdrop-blur-md">
            <div className="mx-auto max-w-7xl px-4 py-4 sm:px-6 lg:px-8">
                <div className="flex flex-col gap-4 lg:flex-row lg:items-center lg:justify-between">
                    {/* 左侧：品牌和标题 */}
                    <div className="flex items-center gap-4">
                        <div className="flex h-10 w-10 items-center justify-center">
                            <img src={'/logo.png'} className="h-10 w-10" alt={'logo'}/>
                        </div>
                        <div>
                            <p className="text-[11px] font-semibold uppercase tracking-[0.4em] text-blue-600">
                                Pika Monitor
                            </p>
                            <h1 className="mt-0.5 text-xl font-bold text-slate-900 lg:text-2xl">{title}</h1>
                        </div>
                    </div>

                    {/* 右侧：操作区域 */}
                    <div className="flex flex-wrap items-center gap-2 text-xs">
                        {/* 最后更新时间 */}
                        {lastUpdated && (
                            <>
                                <div className="flex items-center gap-1.5 rounded-lg bg-slate-50 px-3 py-2 text-slate-600">
                                    <div className="h-2 w-2 animate-pulse rounded-full bg-emerald-500"/>
                                    <span className="text-xs">
                                        最后更新：<span className="font-semibold text-slate-900">{lastUpdated}</span>
                                    </span>
                                </div>
                                <span className="hidden h-4 w-px bg-slate-200 lg:inline-block"/>
                            </>
                        )}

                        {/* 视图切换 */}
                        {showViewToggle && viewMode && onViewModeChange && (
                            <>
                                <div className="inline-flex items-center gap-0.5 rounded-lg border border-slate-200 bg-slate-50 p-0.5 shadow-sm">
                                    <button
                                        type="button"
                                        onClick={() => onViewModeChange('grid')}
                                        className={`inline-flex items-center gap-1 rounded-md px-2 py-1.5 text-xs font-medium transition-all ${
                                            viewMode === 'grid'
                                                ? 'bg-blue-600 text-white shadow-sm'
                                                : 'text-slate-500 hover:bg-white hover:text-blue-600'
                                        }`}
                                        title="网格视图"
                                    >
                                        <LayoutGrid className="h-3.5 w-3.5"/>
                                    </button>
                                    <button
                                        type="button"
                                        onClick={() => onViewModeChange('list')}
                                        className={`inline-flex items-center gap-1 rounded-md px-2 py-1.5 text-xs font-medium transition-all ${
                                            viewMode === 'list'
                                                ? 'bg-blue-600 text-white shadow-sm'
                                                : 'text-slate-500 hover:bg-white hover:text-blue-600'
                                        }`}
                                        title="列表视图"
                                    >
                                        <List className="h-3.5 w-3.5"/>
                                    </button>
                                </div>
                                <span className="hidden h-4 w-px bg-slate-200 lg:inline-block"/>
                            </>
                        )}

                        {/* 导航链接 */}
                        <nav className="flex items-center gap-2">
                            {/* 服务器链接 */}
                            <a
                                href="/"
                                className="group inline-flex items-center gap-1.5 rounded-lg border border-slate-200 bg-white px-3 py-2 text-xs font-medium text-slate-600 shadow-sm transition-all hover:border-blue-300 hover:bg-blue-50 hover:text-blue-700 hover:shadow"
                            >
                                <Server className="h-3.5 w-3.5 transition-transform group-hover:scale-110"/>
                                <span className="hidden sm:inline">服务器</span>
                            </a>

                            {/* 监控链接 */}
                            <a
                                href="/monitors"
                                className="group inline-flex items-center gap-1.5 rounded-lg border border-slate-200 bg-white px-3 py-2 text-xs font-medium text-slate-600 shadow-sm transition-all hover:border-blue-300 hover:bg-blue-50 hover:text-blue-700 hover:shadow"
                            >
                                <Activity className="h-3.5 w-3.5 transition-transform group-hover:scale-110"/>
                                <span className="hidden sm:inline">服务监控</span>
                            </a>
                        </nav>

                        <span className="hidden h-4 w-px bg-slate-200 lg:inline-block"/>

                        {/* GitHub 链接 */}
                        <a
                            href="https://github.com/dushixiang/pika"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center justify-center rounded-lg border border-slate-200 bg-white p-2 shadow-sm transition-all hover:border-slate-300 hover:bg-slate-50 hover:shadow"
                            title="查看 GitHub 仓库"
                        >
                            <img src={GithubSvg} className="h-4 w-4" alt="GitHub"/>
                        </a>

                        {/* 登录按钮 */}
                        <a
                            href="/login"
                            target="_blank"
                            rel="noopener noreferrer"
                            className="inline-flex items-center gap-1.5 rounded-lg border border-blue-200 bg-gradient-to-r px-3 py-2 text-xs font-medium text-white shadow-sm shadow-blue-500/30 transition-all"
                        >
                            <LogIn className="h-3.5 w-3.5"/>
                            <span>登录</span>
                        </a>
                    </div>
                </div>
            </div>
        </header>
    );
};

export default PublicHeader;
