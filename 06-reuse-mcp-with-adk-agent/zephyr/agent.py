import os
from datetime import date

from google.adk.agents import Agent
from google.adk.models.lite_llm import LiteLlm

from google.adk.tools.mcp_tool.mcp_toolset import MCPToolset, StdioServerParameters

from google.genai import types

from .prompts import instructions_root


# INITIALIZE:
os.environ["OPENAI_API_KEY"] = "tada"
os.environ["OPENAI_API_BASE"] = f"{os.environ.get('DMR_BASE_URL')}/engines/llama.cpp/v1"

date_today = date.today()

root_agent = Agent(
    model=LiteLlm(model=f"openai/{os.environ.get('MODEL_RUNNER_CHAT_MODEL')}"),
    name="zephyr_agent",
    generate_content_config=types.GenerateContentConfig(
        temperature=0.0, # More deterministic output
    ),
    description=(
        """
        Zephyr is a dungeon master.
        """
    ),
    instruction= instructions_root(),
    tools=[
        MCPToolset(
            connection_params=StdioServerParameters(
                command='docker',
                args=[
                    "run",
                    "-i",
                    "--rm",
                    "mcp-dd:demo",
                ],
            ),
            # Optional: Filter which tools from the MCP server are exposed
            tool_filter=['choisir_un_personnage_par_son_espece', 'detecter_le_vrai_sujet_du_message_utilisateur']
        )
    ],

)