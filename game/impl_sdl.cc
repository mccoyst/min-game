#include "ui.hpp"
#include "game.hpp"
#include <SDL.h>
#include <SDL_opengl.h>
#include <SDL_image.h>
#include <SDL_ttf.h>
#include <cstdarg>
#include <cassert>

namespace{
extern const char *vshader_src;
extern const char *fshader_src;
GLuint make_buffer(GLenum target, const void *data, GLsizei size);
GLuint make_shader(GLenum type, const char *src);
GLuint make_program(GLuint vshader, GLuint fshader);
}

class SdlUi : public ui::Ui {
	SDL_Surface *win;
	unsigned long tick0;

	GLuint vbuff, ebuff;
	GLuint vshader, fshader, program;
	GLint texloc, posloc, offsloc, shadeloc, dimsloc;
public:
	SdlUi(Fixed w, Fixed h, const char *title);
	~SdlUi();
	virtual void Flip();
	virtual void Clear();
	virtual void Delay(unsigned long);
	virtual unsigned long Ticks();
	virtual bool PollEvent(ui::Event&);
	virtual void Draw(const Vec3&, std::shared_ptr<ui::Img>, float);
};

struct SdlImg : public ui::Img {
	GLuint texId;
	unsigned int w, h;

	SdlImg(SDL_Surface*);
	virtual ~SdlImg();
	virtual unsigned int Width() { return w; }
	virtual unsigned int Height() { return h; }
};

struct SdlFont : public ui::Font {
	TTF_Font *font;
	char r, g, b;

	SdlFont(const char *, int, char, char, char);
	virtual ~SdlFont();
};

SdlUi::SdlUi(Fixed w, Fixed h, const char *title) : Ui(w, h) {
	if (SDL_Init(SDL_INIT_VIDEO) == -1)
		throw Failure("Failed to initialized SDL video");
	tick0 = SDL_GetTicks();

	win = SDL_SetVideoMode(w.whole(), h.whole(), 0, SDL_OPENGL);
	if (!win)
		throw Failure("Failed to set SDL video mode");

	int imgflags = IMG_INIT_PNG;
	if ((IMG_Init(imgflags) & imgflags) != imgflags)
		throw Failure("Failed to initialize png support: %s", IMG_GetError());

	if (TTF_Init() == -1)
		throw Failure("Failed to initialize SDL_ttf: %s", TTF_GetError());

	glEnable(GL_ALPHA_TEST);
	glAlphaFunc(GL_GREATER, 0.5);
	gluOrtho2D(0, w.whole(), 0, h.whole());

	GLfloat vertices[] = {
		0.0f, 0.0f, 0, 1,
		1.0f, 0.0f, 1, 1,
		0.0f, 1.0f, 0, 0,
		1.0f, 1.0f, 1, 0,
	};
	GLushort elements[] = { 0, 1, 2, 3 };

	vbuff = make_buffer(GL_ARRAY_BUFFER, vertices, sizeof(vertices));
	ebuff = make_buffer(GL_ELEMENT_ARRAY_BUFFER, elements, sizeof(elements));

	vshader = make_shader(GL_VERTEX_SHADER, vshader_src);
	if(!vshader)
		throw Failure("Failed to compile vertex shader");
	fshader = make_shader(GL_FRAGMENT_SHADER, fshader_src);
	if(!fshader)
		throw Failure("Failed to compile fragment shader");
	program = make_program(vshader, fshader);
	if(!program)
		throw Failure("Failed to link program");

	texloc = glGetUniformLocation(program, "tex");
	posloc = glGetAttribLocation(program, "position");
	offsloc = glGetUniformLocation(program, "offset");
	shadeloc = glGetUniformLocation(program, "shade");
	dimsloc = glGetUniformLocation(program, "dims");
}

SdlUi::~SdlUi() {
	TTF_Quit();
	IMG_Quit();
	SDL_Quit();
}

void SdlUi::Flip() {
	SDL_GL_SwapBuffers();
}

void SdlUi::Clear() {
	glClear(GL_COLOR_BUFFER_BIT);
}

void SdlUi::Delay(unsigned long msec) {
	SDL_Delay(msec);
}

unsigned long SdlUi::Ticks() {
	return tick0 - SDL_GetTicks();
}

static bool getbutton(SDL_Event &sdle, ui::Event &e) {
	switch (sdle.button.button) {
	case SDL_BUTTON_LEFT:
		e.button = ui::Event::MouseLeft;
		break;
	case SDL_BUTTON_RIGHT:
		e.button = ui::Event::MouseRight;
		break;
	case SDL_BUTTON_MIDDLE:
		e.button = ui::Event::MouseCenter;
		break;
	default:
		return false;
	};
	return true;
}

static bool getkey(SDL_Event &sdle, ui::Event &e) {
	switch (sdle.key.keysym.sym) {
	case SDLK_UP:
		e.button = ui::Event::KeyUpArrow;
		break;
	case SDLK_DOWN:
		e.button = ui::Event::KeyDownArrow;
		break;
	case SDLK_LEFT:
		e.button = ui::Event::KeyLeftArrow;
		break;
	case SDLK_RIGHT:
		e.button = ui::Event::KeyRightArrow;
		break;
	case SDLK_RSHIFT:
		e.button = ui::Event::KeyRShift;
		break;
	case SDLK_LSHIFT:
		e.button = ui::Event::KeyLShift;
		break;
	default:
		if (sdle.key.keysym.sym < 'a' || sdle.key.keysym.sym > 'z')
			return false;
		e.button = sdle.key.keysym.sym;
	}

	return true;
}

bool SdlUi::PollEvent(ui::Event &e) {
	SDL_Event sdle;
	while (SDL_PollEvent(&sdle)) {
		switch (sdle.type) {
		case SDL_QUIT:
			e.type = ui::Event::Closed;
			return true;

		case SDL_MOUSEBUTTONDOWN:
			e.type = ui::Event::MouseDown;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(sdle, e))
				continue;
			return true;

		case SDL_MOUSEBUTTONUP:
			e.type = ui::Event::MouseUp;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(sdle, e))
				continue;
			return true;


		case SDL_MOUSEMOTION:
			e.type = ui::Event::MouseMoved;
			e.x = sdle.motion.x;
			e.y = sdle.motion.y;
			return true;

		case SDL_KEYUP:
			e.type = ui::Event::KeyUp;
			if (!getkey(sdle, e))
				continue;
			return true;

		case SDL_KEYDOWN:
			e.type = ui::Event::KeyDown;
			if (!getkey(sdle, e))
				continue;
			return true;

		default:
			// ignore
			break;
		}
	}
	return false;
}

void SdlUi::Draw(const Vec3 &l, std::shared_ptr<ui::Img> _img, float shade) {
	SdlImg *img = static_cast<SdlImg*>(_img.get());
	if(shade < 0) shade = 0;
	else if(shade > 1) shade = 1;

	glUseProgram(program);
	glActiveTexture(GL_TEXTURE0);
	glBindTexture(GL_TEXTURE_2D,img-> texId);
	glUniform1i(texloc, 0);

	glUniform2f(offsloc, l.x.whole(), height.whole() - l.y.whole());
	glUniform1f(shadeloc, shade);
	glUniform2f(dimsloc, img->w, img->h);

	glBindBuffer(GL_ARRAY_BUFFER, vbuff);
	glVertexAttribPointer(posloc, 4, GL_FLOAT, GL_FALSE, sizeof(GLfloat[4]), 0);
	glEnableVertexAttribArray(posloc);
	glBindBuffer(GL_ELEMENT_ARRAY_BUFFER, ebuff);
	glDrawElements(GL_TRIANGLE_STRIP, 4, GL_UNSIGNED_SHORT, 0);
	glDisableVertexAttribArray(posloc);
}

SdlImg::SdlImg(SDL_Surface *surf) : w(surf->w), h(surf->h) {
	GLint pxSz = surf->format->BytesPerPixel;
	GLenum texFormat = GL_BGRA;
	switch (pxSz) {
	case 4:
		if (surf->format->Rmask == 0xFF)
			texFormat = GL_RGBA;
		break;
	case 3:
		if (surf->format->Rmask == 0xFF)
			texFormat = GL_RGB;
		else
			texFormat = GL_BGR;
		break;
	default:
		throw Failure("Bad image color type… apparently");
	}

	glGenTextures(1, &texId);
	glBindTexture(GL_TEXTURE_2D, texId);
 
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_LINEAR);
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_LINEAR);
 
	glTexImage2D(GL_TEXTURE_2D, 0, pxSz, surf->w, surf->h, 0,
		texFormat, GL_UNSIGNED_BYTE, surf->pixels);
}

SdlImg::~SdlImg() {
	glDeleteTextures(1, &texId);
}

SdlFont::SdlFont(const char *path, int sz, char _r, char _g, char _b)
		: r(_r), g(_g), b(_b) {
	font = TTF_OpenFont(path, sz);
	if (!font)
		throw Failure("Failed to load font %s: %s", path, TTF_GetError());
}

SdlFont::~SdlFont() {
	TTF_CloseFont(font);
}

std::shared_ptr<ui::Ui> ui::OpenWindow(Fixed w, Fixed h, const char *title) {
	return std::shared_ptr<ui::Ui>(new SdlUi(w, h, title));
}

std::shared_ptr<ui::Img> ui::LoadImg(const char *path) {
	SDL_Surface *surf = IMG_Load(path);
	if (!surf)
		throw Failure("Failed to load image %s", path);
	std::shared_ptr<ui::Img> img(new SdlImg(surf));
	SDL_FreeSurface(surf);
	return img;
}

std::shared_ptr<ui::Font> ui::LoadFont(const char *path, int sz, char r, char g, char b) {
	return std::shared_ptr<ui::Font>(new SdlFont(path, sz, r, g, b));
}

std::shared_ptr<ui::Img> ui::RenderText(std::shared_ptr<ui::Font> f, const char *fmt, ...) {
	SdlFont *font = dynamic_cast<SdlFont*>(f.get());
	assert (font);

	char s[256];
	va_list ap;
	va_start(ap, fmt);
	vsnprintf(s, sizeof(s), fmt, ap);
	va_end(ap);

	SDL_Color c;
	c.r = font->r;
	c.g = font->g;
	c.b = font->b;
	SDL_Surface *surf = TTF_RenderUTF8_Blended(font->font, s, c);
	if (!surf)
		throw Failure("Failed to render text: %s", TTF_GetError());

	std::shared_ptr<ui::Img> img(new SdlImg(surf));
	SDL_FreeSurface(surf);
	return img;
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
	glShaderSource(shader, 1, &src, &len);
	glCompileShader(shader);
	glGetShaderiv(shader, GL_COMPILE_STATUS, &shader_ok);
	if(!shader_ok){
		GLint log_len;
		glGetShaderiv(shader, GL_INFO_LOG_LENGTH, &log_len);
		char *log = new char[len];
		glGetShaderInfoLog(shader, log_len, NULL, log);
		fprintf(stderr, "Shader error: %s", log);
		delete [] log;
		glDeleteShader(shader);
		return 0;
	}
	return shader;
}

GLuint make_program(GLuint vshader, GLuint fshader){
	GLint program_ok;

	GLuint program = glCreateProgram();
	glAttachShader(program, vshader);
	glAttachShader(program, fshader);
	glLinkProgram(program);
	glGetProgramiv(program, GL_LINK_STATUS, &program_ok);
	if(!program_ok){
		GLint log_len;
		glGetProgramiv(program, GL_INFO_LOG_LENGTH, &log_len);
		char *log = new char[log_len];
		glGetProgramInfoLog(program, log_len, NULL, log);
		fprintf(stderr, "Program error: %s", log);
		delete [] log;
		glDeleteProgram(program);
		return 0;
	}
	return program;
}

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
}
