#version 120

attribute vec4 position;
attribute int texid;
attribute float shade;

varying vec2 texCoord;
varying int texId;
varying float texShade;

uniform vec2 offset;

void main(){
	vec4 trans = vec4(position.xy + offset, 0.0, 1.0);
	gl_Position = gl_ModelViewProjectionMatrix * trans;

	texcoord = position.zw;
	texId = texid;
	texShade = shade;
}
