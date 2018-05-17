#version 330 core
out vec3 fragCol;
  
in vec2 fragPos;
in vec2 fragUV;

uniform sampler2D tex0;

void main()
{
    fragCol = texture(tex0, vec2(fragUV.x, -fragUV.y)).rgb;
}