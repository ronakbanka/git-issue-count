<html>
<head>
<title>Check Github Issues</title>
<link href="/static/css/main.css" rel="stylesheet" type="text/css">
</head>
<body>
<section class="webdesigntuts-workshop">
<form action="/getinfo" method="post">
    <input type="text" name="repo" placeholder="Github URL">
    <input type="submit" value="Get Issue Count">
</form>
<p>
{{ .Data }}
</p>
</section>
</body>
</html>
