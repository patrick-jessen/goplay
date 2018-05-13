#version 330 core
layout (location = 0) in vec3 vertPos;
layout (location = 1) in vec3 vertNorm;
layout (location = 2) in vec2 vertUV;

layout (std140) uniform shader_data {
  mat4 viewProjMat;
  mat4 modelMat;
};

out vec3 fragNorm;
out vec2 fragUV;

void main() {
  gl_Position = viewProjMat * modelMat * vec4(vertPos, 1.0);
  fragNorm = mat3(modelMat) * vertNorm;
  fragUV = vertUV;
}