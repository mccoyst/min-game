#version 120

varying vec2 texCoord;
varying int texId;
varying float texShade;

uniform sampler2D texes[3];

void main(){
	vec4 c = texture2D(texes[texId], texCoord);
	gl_FragColor = vec4(c.rgb*texShade, tc.a);
}
