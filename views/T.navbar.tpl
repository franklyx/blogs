{{define "navbar"}}
<div class="navbar navbar-default navbar-fixed-top">
    <div class="container">
        <div>
        <a href="/" class="navbar-brand">我的博客</a>
        <ul class="nav navbar-nav">
            <li {{if .IsHome}} class="active" {{end}}><a href="/">首页</a></li>
            <li {{if .IsCategory}} class="active" {{end}}><a href="/category">分类</a></li>
            <li {{if .IsTopic}} class="active" {{end}}><a href="/topic">文章</a></li>
        </ul>
        </div>
        <div class="pull-right">
         <ul class="nav navbar-nav">
            {{if .IsLogin}}
            <li class=""><a href="/login?exit=true">管理员退出 </a></li>
            {{else}}
            <li class=""><a href="/login">管理员登陆</a></li>
            {{end}}
         </ul>
        </div>
    </div>

</div>
{{end}}