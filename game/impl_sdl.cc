#include "ui.hpp"
#include "game.hpp"
#include <SDL.h>
#include <SDL_opengl.h>
#include <SDL_image.h>
#include <SDL_ttf.h>
#include <cstdarg>

namespace{
extern const char *vshader_src;
extern const char *fshader_src;
GLuint make_buffer(GLenum target, const void *data, GLsizei size);
GLuint make_shader(GLenum type, const char *src);
GLuint make_program(GLuint vshader, GLuint fshader);
}

class SdlUi : public Ui {
	SDL_Surface *win;

	GLuint vbuff;
	GLuint vshader, fshader, program;
	GLint texloc, posloc, offsloc, shadeloc, dimsloc;
public:
	SdlUi(Fixed w, Fixed h, const char *title);
	~SdlUi();
	virtual void Flip();
	virtual void Clear();
	virtual void Delay(unsigned long);
	virtual unsigned long Ticks();
	virtual bool PollEvent(Event&);
	virtual void Draw(const Vec2&, std::shared_ptr<Img>, float);
};

struct SdlImg : public Img {
	GLuint texId;
	Vec2 sz;

	SdlImg(SDL_Surface*);
	virtual ~SdlImg();
	virtual Vec2 Size() const { return sz; }
};

struct SdlFont : public Font {
	TTF_Font *font;
	char r, g, b;

	SdlFont(const char *, int, char, char, char);
	virtual ~SdlFont();
	virtual std::shared_ptr<Img> Render(const char*, ...);
};

SdlUi::SdlUi(Fixed w, Fixed h, const char *title) : Ui(w, h) {
	if (SDL_Init(SDL_INIT_VIDEO) == -1)
		throw Failure("Failed to initialized SDL video");

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

	vbuff = make_buffer(GL_ARRAY_BUFFER, vertices, sizeof(vertices));

	vshader = make_shader(GL_VERTEX_SHADER, vshader_src);
	fshader = make_shader(GL_FRAGMENT_SHADER, fshader_src);
	program = make_program(vshader, fshader);

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
	return SDL_GetTicks();
}

static bool getbutton(SDL_Event &sdle, Event &e) {
	switch (sdle.button.button) {
	case SDL_BUTTON_LEFT:
		e.button = Event::MouseLeft;
		break;
	case SDL_BUTTON_RIGHT:
		e.button = Event::MouseRight;
		break;
	case SDL_BUTTON_MIDDLE:
		e.button = Event::MouseCenter;
		break;
	default:
		return false;
	};
	return true;
}

static bool getkey(SDL_Event &sdle, Event &e) {
	switch (sdle.key.keysym.sym) {
	case SDLK_UP:
		e.button = Event::KeyUpArrow;
		break;
	case SDLK_DOWN:
		e.button = Event::KeyDownArrow;
		break;
	case SDLK_LEFT:
		e.button = Event::KeyLeftArrow;
		break;
	case SDLK_RIGHT:
		e.button = Event::KeyRightArrow;
		break;
	case SDLK_RSHIFT:
		e.button = Event::KeyRShift;
		break;
	case SDLK_LSHIFT:
		e.button = Event::KeyLShift;
		break;
	default:
		if (sdle.key.keysym.sym < 'a' || sdle.key.keysym.sym > 'z')
			return false;
		e.button = sdle.key.keysym.sym;
	}

	return true;
}

bool SdlUi::PollEvent(Event &e) {
	SDL_Event sdle;
	while (SDL_PollEvent(&sdle)) {
		switch (sdle.type) {
		case SDL_QUIT:
			e.type = Event::Closed;
			return true;

		case SDL_MOUSEBUTTONDOWN:
			e.type = Event::MouseDown;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(sdle, e))
				continue;
			return true;

		case SDL_MOUSEBUTTONUP:
			e.type = Event::MouseUp;
			e.x = sdle.button.x;
			e.y = sdle.button.y;
			if (!getbutton(sdle, e))
				continue;
			return true;


		case SDL_MOUSEMOTION:
			e.type = Event::MouseMoved;
			e.x = sdle.motion.x;
			e.y = sdle.motion.y;
			return true;

		case SDL_KEYUP:
			e.type = Event::KeyUp;
			if (!getkey(sdle, e))
				continue;
			return true;

		case SDL_KEYDOWN:
			e.type = Event::KeyDown;
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

void SdlUi::Draw(const Vec2 &l, std::shared_ptr<Img> _img, float shade) {
	SdlImg *img = static_cast<SdlImg*>(_img.get());
	if(shade < 0) shade = 0;
	else if(shade > 1) shade = 1;

	glUseProgram(program);
	glActiveTexture(GL_TEXTURE0);
	glBindTexture(GL_TEXTURE_2D,img-> texId);
	glUniform1i(texloc, 0);

	glUniform2f(offsloc, l.x.whole(), l.y.whole());
	glUniform1f(shadeloc, shade);
	glUniform2f(dimsloc, img->sz.x.whole(), img->sz.y.whole());

	glBindBuffer(GL_ARRAY_BUFFER, vbuff);
	glVertexAttribPointer(posloc, 4, GL_FLOAT, GL_FALSE, 0, 0);
	glEnableVertexAttribArray(posloc);
	glDrawArrays(GL_TRIANGLE_STRIP, 0, 4);
	glDisableVertexAttribArray(posloc);
}

SdlImg::SdlImg(SDL_Surface *surf) : sz(Fixed(surf->w), Fixed(surf->h)) {
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
 
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MIN_FILTER, GL_NEAREST);
	glTexParameteri(GL_TEXTURE_2D, GL_TEXTURE_MAG_FILTER, GL_NEAREST);
 
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

std::shared_ptr<Img> SdlFont::Render(const char *fmt, ...) {
	char s[256];
	va_list ap;
	va_start(ap, fmt);
	vsnprintf(s, sizeof(s), fmt, ap);
	va_end(ap);

	SDL_Color c;
	c.r = r;
	c.g = g;
	c.b = b;
	SDL_Surface *surf = TTF_RenderUTF8_Blended(font, s, c);
	if (!surf)
		throw Failure("Failed to render text: %s", TTF_GetError());

	std::shared_ptr<Img> img(new SdlImg(surf));
	SDL_FreeSurface(surf);
	return img;
}

std::shared_ptr<Ui> OpenWindow(Fixed w, Fixed h, const char *title) {
	return std::shared_ptr<Ui>(new SdlUi(w, h, title));
}

std::shared_ptr<Img> LoadImg(const char *path) {
	SDL_Surface *surf = IMG_Load(path);
	if (!surf)
		throw Failure("Failed to load image %s", path);
	std::shared_ptr<Img> img(new SdlImg(surf));
	SDL_FreeSurface(surf);
	return img;
}

std::shared_ptr<Font> LoadFont(const char *path, int sz, char r, char g, char b) {
	return std::shared_ptr<Font>(new SdlFont(path, sz, r, g, b));
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
		abort();
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
		abort();
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
