#version 330 core
layout (location = 0) in vec3 vertPos;
layout (location = 2) in vec2 vertUV;

out vec2 fragPos;
out vec2 fragUV;

void main()
{
    gl_Position = vec4(vertPos.xy, 0, 1);
    fragPos = vertPos.xy;
    fragUV = vertUV;
}