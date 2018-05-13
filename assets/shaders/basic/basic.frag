#version 330 core
layout (location = 0) out vec3 fragColor;

in vec3 fragNorm;
in vec2 fragUV;

uniform sampler2D tex0;

vec3 fixedLight = normalize(vec3(4, 1, -1));
float ambient = 0.3;


void main() {
  vec3 light = vec3(1, 1, 1) * (max(dot(fragNorm, fixedLight), ambient)); 
  fragColor = texture(tex0, fragUV).rgb * light;
}