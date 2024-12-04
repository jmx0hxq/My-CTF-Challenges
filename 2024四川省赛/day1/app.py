from flask import Flask, render_template, request, redirect, url_for
from flask_mako import MakoTemplates, render_template as mako_render_template
from mako.template import Template as Mako_T
# from flask_mako import MakoTemplates

app = Flask(__name__)
mako = MakoTemplates(app)

welcome_string = """
<!DOCTYPE html>
<html>
<head>
    <title>My APP</title>
    <style>
        body {
            font-family: Arial, sans-serif;
        }
        .header {
            background-color: #f2f2f2;
            padding: 10px;
            text-align: left;
        }
        .body {
            padding: 20px;
        }
    </style>
</head>
<body>
    <div class="header"> 
   Welcome %s !
   %%if title:
    This your admin page.
   %%endif
    </div>
    <div class="body">
        <p>This is your profile。</p>
    </div>
</body>
</html>
"""
welcome_message = "welcome"
black_list = ['${', 'import', 'os', 'system', 'popen', 'join', 'context', 'sys', '__', 'builtins', 'eval', 'exec', 'ord']

@app.route('/', methods=['GET', 'POST'])
def index():
    if request.method == 'POST':
        username = request.form['username']
        return redirect(url_for('welcome', username=username))
    return mako_render_template('index.html')

@app.route('/welcome')
def welcome():
    username = request.args.get('username')
    if len(username) > 42:
        return "error username"
    for key in black_list:
        if key.upper() in username.upper():
            return "bad username"
    if username == "Admin":
        return Mako_T(welcome_string % username).render(title=True)
    return Mako_T(welcome_string % username).render(title=False)


if __name__ == '__main__':
    app.run(host='0.0.0.0', port=5002)
