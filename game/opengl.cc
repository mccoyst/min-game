#include "opengl.hpp"
#include "game.hpp"
#include <SDL_opengl.h>
#include <cstring>

namespace{
extern const char *vshader_src;
extern const char *fshader_src;
extern const char *world_vshader;
extern const char *world_fshader;
GLuint make_buffer(GLenum target, const void *data, GLsizei size);
GLuint make_shader(GLenum type, const char *src);
GLuint make_program(GLuint vshader, GLuint fshader);
}

void OpenGLUi::InitOpenGL() {
	glEnable(GL_TEXTURE_2D);
	glEnable(GL_BLEND);
	glBlendFunc(GL_SRC_ALPHA, GL_ONE_MINUS_SRC_ALPHA);
	glMatrixMode(GL_PROJECTION);
	glLoadIdentity();
	glClearColor(0.0, 0.0, 0.0, 0.0);
	gluOrtho2D(0, width.whole(), 0, height.whole());
	glMatrixMode(GL_MODELVIEW);
	glLoadIdentity();
}

void OpenGLUi::Clear() {
	glClear(GL_COLOR_BUFFER_BIT);
}

void OpenGLUi::Draw(const Vec2 &l, std::shared_ptr<Img> _img, float shade) {
	OpenGLImg *img = static_cast<OpenGLImg*>(_img.get());
	float x = l.x.whole(), y = l.y.whole();
	float w = img->Size().x.whole(), h = img->Size().y.whole();

	glBindTexture(GL_TEXTURE_2D, img->texid);

	if (shade < 0)
		shade = 0;
	else if (shade > 1)
		shade = 1;

	glColor4f(shade, shade, shade, 1);
	glBegin(GL_QUADS);
	glTexCoord2i(0, 1);
	glVertex3f(x, y, 0);
	glTexCoord2i(1, 1);
	glVertex3f(x+w, y, 0);
	glTexCoord2i(1, 0);
	glVertex3f(x+w, y+h, 0);
	glTexCoord2i(0, 0);
	glVertex3f(x, y+h, 0);
	glEnd();
}

OpenGLImg::~OpenGLImg() {
	glDeleteTextures(1, &texid);
}

namespace{
GLuint make_buffer(GLenum target, const void *data, GLsizei size){
	GLuint buffer;
	glGenBuffers(1, &buffer);
	glBindBuffer(target, buffer);
	glBufferData(target, size, data, GL_STATIC_DRAW);
	return buffer;
}

GLuint make_shader(GLenum type, const char *src){
	GLint len = strlen(src);
	GLuint shader;
	GLint shader_ok;

	shader = glCreateShader(type);
	if (!shader)
		throw Failure("Failed to create a shader");
	glShaderSource(shader, 1, &src, &len);
	glCompileShader(shader);
	glGetShaderiv(shader, GL_COMPILE_STATUS, &shader_ok);
	if(!shader_ok){
		GLint log_len = 0;
		glGetShaderiv(shader, GL_INFO_LOG_LENGTH, &log_len);
		char *log;
		if (log_len > 0) {
			log = new char[log_len];
			glGetShaderInfoLog(shader, log_len, NULL, log);
		} else {
			log = (char*) "<no message>";
		}
		Failure fail("Failed to compile shader: %s", log);
		glDeleteShader(shader);
		if (log_len > 0)
			delete [] log;
		throw fail;
	}
	return shader;
}

GLuint make_program(GLuint vshader, GLuint fshader){
	GLint program_ok;

	GLuint program = glCreateProgram();
	if (!program)
		throw Failure("Failed to create a program");
	glAttachShader(program, vshader);
	glAttachShader(program, fshader);
	glLinkProgram(program);
	glGetProgramiv(program, GL_LINK_STATUS, &program_ok);
	if(!program_ok){
		GLint log_len;
		glGetProgramiv(program, GL_INFO_LOG_LENGTH, &log_len);
		char *log;
		if (log_len > 0) {
			log = new char[log_len];
			glGetProgramInfoLog(program, log_len, NULL, log);
		} else {
			log = (char*) "<no message>";
		}
		Failure fail("Failed to link program: %s", log);
		if (log_len > 0)
			delete [] log;
		glDeleteProgram(program);
		throw fail;
	}
	return program;
}

/*
const char *vshader_src = 
	"#version 120\n"
	"attribute vec4 position;"
	"varying vec2 texcoord;"
	"uniform vec2 offset;"
	"uniform vec2 dims;"
	""
	"void main()"
	"{"
	"	vec2 p = vec2(position.x*dims.x, position.y*dims.y);"
	"	vec4 trans = vec4(p+offset, 0.0, 1.0);"
	"	gl_Position = gl_ModelViewProjectionMatrix * trans;"
	"	texcoord = position.ba;"
	"}"
	;

const char *fshader_src =
	"#version 120\n"
	"uniform sampler2D tex;"
	"uniform float shade;"
	"varying vec2 texcoord;"
	
	"void main()"
	"{"
		"vec4 tc = texture2D(tex, texcoord);"
	"	gl_FragColor = vec4(tc.rgb*shade, tc.a);"
	"}"
	;

const char *world_vshader =
	"#version 120\n"
	""
	"attribute vec4 position;"
	"attribute float texid;"
	"attribute float shade;"
	""
	"varying vec2 texCoord;"
	"varying float texId;"
	"varying float texShade;"
	""
	"uniform vec2 offset;"
	""
	"void main(){"
	"	vec4 trans = vec4(position.xy + offset, 0.0, 1.0);"
	"	gl_Position = gl_ModelViewProjectionMatrix * trans;"
	""
	"	texCoord = position.zw;"
	"	texId = texid;"
	"	texShade = shade;"
	"}"
	;

const char *world_fshader =
	"#version 120\n"
	""
	"varying vec2 texCoord;"
	"varying float texId;"
	"varying float texShade;"
	""
	"uniform sampler2D texes[3];"
	""
	"void main(){"
	"	vec4 c = texture2D(texes[int(texId)], texCoord);"
	"	gl_FragColor = vec4(c.rgb*texShade, c.a);"
	"}"
	;

*/
}