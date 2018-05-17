#version 330 core
layout (location = 0) in vec3 vertPos;
layout (location = 1) in vec3 vertNorm;
layout (location = 2) in vec2 vertUV;
layout (location = 3) in vec3 vertTang;

layout (std140) uniform shader_data {
  mat4 viewProjMat;
  mat4 modelMat;
  vec4 viewPos;
};

out vec3 fragPos;
out vec2 fragUV;
out mat3 TBN;

void main() {
  gl_Position = viewProjMat * modelMat * vec4(vertPos, 1.0);
  fragPos = vec3(modelMat * vec4(vertPos, 1.0));
  fragUV = vertUV;

  vec3 biTan = cross(vertNorm, vertTang);
  vec3 T = normalize(vec3(modelMat * vec4(vertTang, 0.0)));
  vec3 B = normalize(vec3(modelMat * vec4(biTan, 0.0)));
  vec3 N = normalize(vec3(modelMat * vec4(vertNorm, 0.0)));
  TBN = mat3(T, B, N);
}